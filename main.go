package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukhttp"
	"github.com/ipthomas/tukutil"
	"github.com/ipthomas/tukxdw"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var dbconn = tukdbint.TukDBConnection{DBUser: os.Getenv(tukcnst.ENV_DB_USER), DBPassword: os.Getenv(tukcnst.ENV_DB_PASSWORD), DBHost: os.Getenv(tukcnst.ENV_DB_HOST), DBPort: os.Getenv(tukcnst.ENV_DB_PORT), DBName: os.Getenv(tukcnst.ENV_DB_NAME)}
var initstate bool
var deBugMode = false

type def struct {
	Idmaps []tukdbint.IdMap
}

func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.SetFlags(log.Lshortfile)
	var err error
	deBugMode, _ = strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	log.Printf("Processing API Gateway %s Request Path %s", req.HTTPMethod, req.Path)

	if !initstate {
		if err = tukdbint.NewDBEvent(&dbconn); err != nil {
			return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
		}
		initstate = true
	}
	switch req.HTTPMethod {

	case http.MethodGet:
		act := req.QueryStringParameters[tukcnst.ACT]
		//var buf bytes.Buffer
		var data = make(map[string]interface{})
		var dataBytes []byte
		data["XDW_Consumer"] = os.Getenv(tukcnst.ENV_XDW_CONSUMER_URL)
		data["XDW_Creator"] = os.Getenv(tukcnst.ENV_XDW_CREATOR_URL)
		data["XDW_Admin"] = os.Getenv(tukcnst.ENV_XDW_ADMIN_URL)
		data["EVENT_Subscriber"] = os.Getenv(tukcnst.ENV_EVENT_SUBSCRIBER_URL)
		data["User"] = req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_USER]
		data["Org"] = req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_ORG]
		data["Role"] = req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_ROLE]
		data["Act"] = req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_ACT]
		data["Email"] = req.QueryStringParameters["email"]
		data["Task"] = req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_EXPRESSION]

		if dataBytes, err = json.MarshalIndent(data, "", "  "); err != nil {
			return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
		}
		awsReq := tukhttp.AWS_APIRequest{}
		switch act {
		case tukcnst.XDW_ADMIN_REGISTER_DEFINITION, tukcnst.XDW_ADMIN_Register_Template:
			awsReq = tukhttp.AWS_APIRequest{URL: os.Getenv(tukcnst.ENV_HTML_CREATOR_URL), Resource: "uploadfile", Body: dataBytes}
		case "codemap":
			return getCodemap(req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_FORMAT])
		default:
			awsReq = tukhttp.AWS_APIRequest{URL: os.Getenv(tukcnst.ENV_HTML_CREATOR_URL), Resource: "spa", Body: dataBytes}
		}
		if err = tukhttp.NewRequest(&awsReq); err != nil {
			return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
		}
		return queryResponse(http.StatusOK, string(awsReq.Response), tukcnst.TEXT_HTML)

	case http.MethodPost:
		if req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_ACT] == "codemap" {
			return updateCodemap(req)
		}
		var fileBytes []byte
		var form *multipart.Form
		var file multipart.File
		body := req.Body
		contentType := req.Headers["content-type"]
		boundary := strings.Split(contentType, "boundary=")[1]
		reader := multipart.NewReader(bytes.NewReader([]byte(body)), boundary)
		form, err = reader.ReadForm(32 << 20) // 32 MB maximum request size
		if err != nil {
			return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
		}

		fileHeader := form.File["file"][0]
		file, err = fileHeader.Open()
		if err != nil {
			return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
		}
		defer file.Close()
		log.Printf("File name: %s\n", fileHeader.Filename)
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
		}
		log.Printf("File size: %d bytes\n", len(fileBytes))
		act := form.Value["act"][0]
		switch act {
		case tukcnst.XDW_ADMIN_REGISTER_DEFINITION:

			wfdef := tukxdw.WorkflowDefinition{}
			wfmeta := tukxdw.XDSDocumentMeta{}
			pwy := ""
			ismeta := false
			if err = json.Unmarshal(fileBytes, &wfdef); err != nil {
				return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
			}
			pwy = wfdef.Ref
			if pwy == "" {
				if err = json.Unmarshal(fileBytes, &wfmeta); err != nil {
					return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
				}
				pwy = strings.TrimSuffix(wfmeta.ID, "_meta")
				ismeta = true
			}
			if pwy == "" {
				log.Println("unable to set pathway from input file id / name")
				return queryResponse(http.StatusBadRequest, "unable to set pathway from input file id / name", tukcnst.TEXT_PLAIN)
			}
			log.Printf("Pathway value: %s Is Meta %v", pwy, ismeta)
			// fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)
			urls := tukxdw.ServiceURL{DSUB_Broker: os.Getenv("DSUB_BROKER_URL"), DSUB_Consumer: os.Getenv("DSUB_CONSUMER_URL")}
			trans := tukxdw.Transaction{ServiceURL: urls, Request: fileBytes, Pathway: pwy}
			err = trans.RegisterWorkflowDefinition(ismeta)
			if err != nil {
				return queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
			}
			return queryResponse(http.StatusOK, string(fileBytes), tukcnst.APPLICATION_JSON)
		case tukcnst.XDW_ADMIN_Register_Template:
			if strings.HasSuffix(fileHeader.Filename, ".html") {
				if fileBytes != nil {
					log.Printf("Persisting HTML Template %s", fileHeader.Filename)
					tmplts := tukdbint.Templates{Action: tukcnst.DELETE}
					tmplt := tukdbint.Template{Name: strings.TrimSuffix(fileHeader.Filename, ".html")}
					tmplts.Templates = append(tmplts.Templates, tmplt)
					tukdbint.NewDBEvent(&tmplts)
					tmplts = tukdbint.Templates{Action: tukcnst.INSERT}
					tmplt = tukdbint.Template{Name: strings.TrimSuffix(fileHeader.Filename, ".html"), Template: string(fileBytes)}
					tmplts.Templates = append(tmplts.Templates, tmplt)
					tukdbint.NewDBEvent(&tmplts)
					return queryResponse(http.StatusOK, string(fileBytes), tukcnst.TEXT_PLAIN)
				}
			}
		}
	}
	return queryResponse(http.StatusBadRequest, "", tukcnst.TEXT_PLAIN)
}
func updateCodemap(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Println("act = codemap")
	var formvals url.Values
	formvals, _ = url.ParseQuery(req.Body)
	if deBugMode {
		for k, v := range formvals {
			log.Printf("form key %s form val %s", k, v[0])
		}
	}
	idmaps := tukdbint.IdMaps{Action: tukcnst.SELECT}
	idmap := tukdbint.IdMap{Id: int64(tukutil.GetIntFromString(formvals.Get("id")))}
	idmaps.LidMap = append(idmaps.LidMap, idmap)
	tukdbint.NewDBEvent(&idmaps)

	if formvals.Get(tukcnst.QUERY_PARAM_ACTION) == "edit" {
		newidmaps := tukdbint.IdMaps{Action: tukcnst.UPDATE}
		newidmap := tukdbint.IdMap{Id: idmap.Id, Lid: idmaps.LidMap[1].Lid, Mid: idmaps.LidMap[1].Mid}

		if formvals.Has("lid") {
			newidmap.Lid = formvals.Get("lid")
		}
		if formvals.Has("mid") {
			newidmap.Mid = formvals.Get("mid")
		}
		newidmaps.LidMap = append(newidmaps.LidMap, newidmap)
		tukdbint.NewDBEvent(&newidmaps)
	}
	if formvals.Get(tukcnst.QUERY_PARAM_ACTION) == "newCodeMap" {
		newidmaps := tukdbint.IdMaps{Action: tukcnst.INSERT}
		newidmap := tukdbint.IdMap{}
		if formvals.Has("lid") && formvals.Has("mid") {
			newidmap.Lid = formvals.Get("lid")
			newidmap.Mid = formvals.Get("mid")
			newidmaps.LidMap = append(newidmaps.LidMap, newidmap)
			tukdbint.NewDBEvent(&newidmaps)
			log.Printf("Created New Codemap with Local ID %s and Mapped ID %s", newidmap.Lid, newidmap.Mid)
		}
	}
	return getCodemap(tukcnst.TEXT_HTML)
}
func getCodemap(format string) (*events.APIGatewayProxyResponse, error) {
	idmaps := tukdbint.IdMaps{Action: tukcnst.SELECT}
	if err := tukdbint.NewDBEvent(&idmaps); err != nil {
		log.Println(err.Error())
	}

	if format == "json" {
		rsp, _ := json.Marshal(idmaps.LidMap)
		return queryResponse(http.StatusOK, string(rsp), tukcnst.APPLICATION_JSON)
	}

	def := def{Idmaps: idmaps.LidMap}
	sort.Sort(def)
	idmaps.LidMap = def.Idmaps
	rsp, _ := json.Marshal(idmaps)

	awsReq := tukhttp.AWS_APIRequest{URL: os.Getenv(tukcnst.ENV_HTML_CREATOR_URL), Act: tukcnst.SELECT, Resource: "codemaps", Body: rsp}
	tukhttp.NewRequest(&awsReq)
	return queryResponse(http.StatusOK, string(awsReq.Response), tukcnst.TEXT_HTML)
}
func (e def) Len() int {
	return len(e.Idmaps)
}
func (e def) Less(i, j int) bool {
	return e.Idmaps[i].Lid < e.Idmaps[j].Lid
}
func (e def) Swap(i, j int) {
	e.Idmaps[i], e.Idmaps[j] = e.Idmaps[j], e.Idmaps[i]
}
func queryResponse(statusCode int, body string, contentType string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    setAwsResponseHeaders(contentType),
		Body:       body,
	}, nil
}
func setAwsResponseHeaders(contentType string) map[string]string {
	awsHeaders := make(map[string]string)
	awsHeaders["Server"] = "XDW_Admin"
	awsHeaders["Access-Control-Allow-Origin"] = "*"
	awsHeaders["Access-Control-Allow-Headers"] = "accept, Content-Type"
	awsHeaders["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
	awsHeaders[tukcnst.CONTENT_TYPE] = contentType
	return awsHeaders
}
