# tukpdq
tukpdq provides a golang implementtion of IHE PIXm, IHE PIXv3 and IHE PDQv3 Consumer clients

There is currently no authentication implemented. The github.com/ipthomas/tukhttp package is used to handle the http request/response and should be amended according to your authentication requirements

Struct PDQQuery implements the tukpdq.PDQ interface

	type PDQQuery struct {
		Server       string
		Server_URL   string
		NHS_ID       string
		NHS_OID      string
		MRN_ID       string
		MRN_OID      string
		REG_ID       string
		REG_OID      string
		Timeout      int64
		Cache        bool
		RspType      string
		Used_PID     string
		Used_PID_OID string
		Request      []byte
		Response     []byte
		StatusCode   int
		Count        int
		Patients     []PIXPatient
	}
	Server must be set to either "pixm" to perform a IHE PIXm query or "pixv3" to perform an IHE PIXv3 query or "pdqv3" to perform an IHE PDQv3 query. The github.com/ipthomas/tukcnst provides constants for each of the valid Server values i.e. tukcnst.PIXm, tukcnst.PIXv3, tukcnst.PDQv3, or you can just use strings!
	
	
	A patient identifier is required for use in the PDQ. Only one id needs to be provided, not all id's are needed!!
		i.e. This can be either the MRN id along with the associated MRN OID or the NHS ID or the XDS regional ID. The default nhs oid will be used if not provided. The regional oid is always reguired even if not using the reg id in the pdq because when parsing the pdq response the reg oid is needed to identify the patient reg id.

	 The REG_OID is the Regional/XDS OID and is required
	 
	 Server_URL is the IHE (PIXm or PIXv3 or PDQv3) compliant server WS end point and is required.

	 Timeout is the http context timeout in seconds and is optional. Default is 5 secs

	 Cache = "true" enables the caching of found patients for the lifetime of the lambda function. Default is false

	 RspType sets the response type sent to the PDQ. 
	 	If set to "bool" the response will be either true or false.
		If set to "code" the response will be empty and the StatusCode will be either 200 if patient exists or 204 if not
		Default is the reponse is empty if patient not found or the patient details if patient is found

	 Used_ID will be set to the ID used for the query

	 Used_OID will be set to the OID used for the query
	 
	 Request will be set with the bytes of the request message

	 Response will be set as determined by the RspType param. Descibed above 

	 StatusCode will be set according to the RspType param.
	 
	 Count will equal the number of patients found matching the query. Hopefully just 1 !!

	 []Patients will contain an array of PIXPatient structs containing all matched patients. Hopefully just 1 !!

	type PIXPatient struct {
		PIDOID     string `json:"pidoid"`
		PID        string `json:"pid"`
		REGOID     string `json:"regoid"`
		REGID      string `json:"regid"`
		NHSOID     string `json:"nhsoid"`
		NHSID      string `json:"nhsid"`
		GivenName  string `json:"givenname"`
		FamilyName string `json:"familyname"`
		Gender     string `json:"gender"`
		BirthDate  string `json:"birthdate"`
		Street     string `json:"street"`
		Town       string `json:"town"`
		City       string `json:"city"`
		State      string `json:"state"`
		Country    string `json:"country"`
		Zip        string `json:"zip"`
	}

	Example usage as AWS Lambda function:
		pdq := tukpdq.PDQQuery{
			Server:     os.Getenv(tukcnst.AWS_ENV_PDQ_SERVER_TYPE),
			MRN_ID:     req.QueryStringParameters[tukcnst.QUERY_PARAM_MRN_ID],
			MRN_OID:    req.QueryStringParameters[tukcnst.QUERY_PARAM_MRN_OID],
			NHS_ID:     req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_ID],
			NHS_OID:    os.Getenv(tukcnst.AWS_ENV_REG_OID),
			REG_ID:     req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_ID],
			REG_OID:    os.Getenv(tukcnst.AWS_ENV_REG_OID),
			Server_URL: os.Getenv(tukcnst.AWS_ENV_PDQ_SERVER_URL),
			Timeout:    2,
		}
	err = tukpdq.New_Transaction(&pdq)

	Running the above example produces the following Log output:

	2022/09/12 14:02:55.510679 tukpdq.go:188: HTTP GET Request Headers

	2022/09/12 14:02:55.510834 tukpdq.go:190: {
	  "Accept": [
	    "*/*"
	  ],
	  "Connection": [
	    "keep-alive"
	  ],
	  "Content-Type": [
	    "application/json"
	  ]
	}

2022/09/12 14:02:55.510860 tukpdq.go:191: HTTP Request
URL = http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?identifier=2.16.840.1.113883.2.1.4.1%7C9999999468&_format=json&_pretty=true
2022/09/12 14:02:55.851605 tukpdq.go:194: HTML Response - Status Code = 200

	{
	  "resourceType": "Bundle",
	  "id": "53c44d32-fb2c-4dfb-b819-db2150e6fa87",
	  "type": "searchset",
	  "total": 1,
	  "link": [ {
	    "relation": "self",
	    "url": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?_format=json&_pretty=true&identifier=2.16.840.1.113883.2.1.4.1%7C9999999468"
	  } ],
	  "entry": [ {
	    "fullUrl": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient/VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
	    "resource": {
	      "resourceType": "Patient",
	      "id": "VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
	      "extension": [ {
	        "url": "http://hl7.org/fhir/StructureDefinition/patient-citizenship",
	        "valueCodeableConcept": {
	          "coding": [ {
	            "code": "GBR"
	          } ]
	        }
	      }, {
	        "url": "http://hl7.org/fhir/StructureDefinition/patient-nationality",
	        "valueCodeableConcept": {
	          "coding": [ {
	            "code": "GBR"
	          } ]
	        }
	      } ],
	      "identifier": [ {
	        "system": "urn:oid:2.16.840.1.113883.2.1.4.1",
	        "value": "9999999468"
	      }, {
	        "use": "usual",
	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1",
	        "value": "TSUK.16619762302611"
	      }, {
	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1",
	        "value": "REG.1MWU5C92M2"
	      } ],
	      "active": true,
	      "name": [ {
	        "use": "official",
	        "family": "Testpatient",
	        "given": [ "Nhs" ]
	      } ],
	      "telecom": [ {
	        "system": "phone",
	        "value": "07777661324",
	        "use": "work"
	      }, {
	        "system": "email",
	        "value": "nhs.testpatient@nhs.net",
	        "use": "work"
	      } ],
	      "gender": "male",
	      "birthDate": "1962-04-04",
	      "address": [ {
	        "line": [ "Preston Road" ],
	        "city": "Preston",
	        "state": "Lancashire",
	        "postalCode": "PR1 1PR",
	        "country": "GBR"
	      } ],
	      "maritalStatus": {
	        "coding": [ {
	          "code": "D"
	        } ]
	      },
	      "multipleBirthBoolean": false
	    }
	  } ]
	}

	2022/09/12 14:02:55.852334 tukpdq.go:102: 1 Patient Entries in Response
	2022/09/12 14:02:55.852392 tukpdq.go:122: Set NHS ID 9999999468 2.16.840.1.113883.2.1.4.1
	2022/09/12 14:02:55.852427 tukpdq.go:117: Set PID TSUK.16619762302611 2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1
	2022/09/12 14:02:55.852455 tukpdq.go:112: Set Reg ID REG.1MWU5C92M2 2.16.840.1.113883.2.1.3.31.2.1.1
	2022/09/12 14:02:55.852546 tukpdq.go:149: Added Patient 9999999468 to response
	2022/09/12 14:02:55.852569 main.go:84: Patient Nhs Testpatient is registered