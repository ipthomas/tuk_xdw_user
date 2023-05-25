package tukdsub

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"strings"
	"text/template"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukhttp"
	"github.com/ipthomas/tukpdq"
	"github.com/ipthomas/tukutil"
)

var DebugMode = false

// DSUBEvent implements NewEvent(i DSUB_Interface) error
type DSUBEvent struct {
	Action          string
	BrokerURL       string
	ConsumerURL     string
	PDQ_SERVER_URL  string
	PDQ_SERVER_TYPE string
	REG_OID         string
	NHS_OID         string
	EventMessage    string
	Expressions     []string
	Pathway         string
	RowID           int64
	DBConnection    tukdbint.TukDBConnection
	Subs            tukdbint.Subscriptions
	Request         []byte
	Response        []byte
	BrokerRef       string
	UUID            string
	Notify          DSUBNotifyMessage
	Event           tukdbint.Event
}

// DSUBSubscribeResponse is an IHE DSUB Subscribe Message compliant struct
type DSUBSubscribeResponse struct {
	XMLName        xml.Name `xml:"Envelope"`
	Text           string   `xml:",chardata"`
	S              string   `xml:"s,attr"`
	A              string   `xml:"a,attr"`
	Xsi            string   `xml:"xsi,attr"`
	Wsnt           string   `xml:"wsnt,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Header         struct {
		Text   string `xml:",chardata"`
		Action string `xml:"Action"`
	} `xml:"Header"`
	Body struct {
		Text              string `xml:",chardata"`
		SubscribeResponse struct {
			Text                  string `xml:",chardata"`
			SubscriptionReference struct {
				Text    string `xml:",chardata"`
				Address string `xml:"Address"`
			} `xml:"SubscriptionReference"`
		} `xml:"SubscribeResponse"`
	} `xml:"Body"`
}

// DSUBNotifyMessage is an IHE DSUB Notify Message compliant struct
type DSUBNotifyMessage struct {
	XMLName             xml.Name `xml:"Notify"`
	Text                string   `xml:",chardata"`
	Xmlns               string   `xml:"xmlns,attr"`
	Xsd                 string   `xml:"xsd,attr"`
	Xsi                 string   `xml:"xsi,attr"`
	NotificationMessage struct {
		Text                  string `xml:",chardata"`
		SubscriptionReference struct {
			Text    string `xml:",chardata"`
			Address struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
			} `xml:"Address"`
		} `xml:"SubscriptionReference"`
		Topic struct {
			Text    string `xml:",chardata"`
			Dialect string `xml:"Dialect,attr"`
		} `xml:"Topic"`
		ProducerReference struct {
			Text    string `xml:",chardata"`
			Address struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
			} `xml:"Address"`
		} `xml:"ProducerReference"`
		Message struct {
			Text                 string `xml:",chardata"`
			SubmitObjectsRequest struct {
				Text               string `xml:",chardata"`
				Lcm                string `xml:"lcm,attr"`
				RegistryObjectList struct {
					Text            string `xml:",chardata"`
					Rim             string `xml:"rim,attr"`
					ExtrinsicObject struct {
						Text       string `xml:",chardata"`
						A          string `xml:"a,attr"`
						ID         string `xml:"id,attr"`
						MimeType   string `xml:"mimeType,attr"`
						ObjectType string `xml:"objectType,attr"`
						Slot       []struct {
							Text      string `xml:",chardata"`
							Name      string `xml:"name,attr"`
							ValueList struct {
								Text  string   `xml:",chardata"`
								Value []string `xml:"Value"`
							} `xml:"ValueList"`
						} `xml:"Slot"`
						Name struct {
							Text            string `xml:",chardata"`
							LocalizedString struct {
								Text  string `xml:",chardata"`
								Value string `xml:"value,attr"`
							} `xml:"LocalizedString"`
						} `xml:"Name"`
						Description    string `xml:"Description"`
						Classification []struct {
							Text                 string `xml:",chardata"`
							ClassificationScheme string `xml:"classificationScheme,attr"`
							ClassifiedObject     string `xml:"classifiedObject,attr"`
							ID                   string `xml:"id,attr"`
							NodeRepresentation   string `xml:"nodeRepresentation,attr"`
							ObjectType           string `xml:"objectType,attr"`
							Slot                 []struct {
								Text      string `xml:",chardata"`
								Name      string `xml:"name,attr"`
								ValueList struct {
									Text  string   `xml:",chardata"`
									Value []string `xml:"Value"`
								} `xml:"ValueList"`
							} `xml:"Slot"`
							Name struct {
								Text            string `xml:",chardata"`
								LocalizedString struct {
									Text  string `xml:",chardata"`
									Value string `xml:"value,attr"`
								} `xml:"LocalizedString"`
							} `xml:"Name"`
						} `xml:"Classification"`
						ExternalIdentifier []struct {
							Text                 string `xml:",chardata"`
							ID                   string `xml:"id,attr"`
							IdentificationScheme string `xml:"identificationScheme,attr"`
							ObjectType           string `xml:"objectType,attr"`
							RegistryObject       string `xml:"registryObject,attr"`
							Value                string `xml:"value,attr"`
							Name                 struct {
								Text            string `xml:",chardata"`
								LocalizedString struct {
									Text  string `xml:",chardata"`
									Value string `xml:"value,attr"`
								} `xml:"LocalizedString"`
							} `xml:"Name"`
						} `xml:"ExternalIdentifier"`
					} `xml:"ExtrinsicObject"`
				} `xml:"RegistryObjectList"`
			} `xml:"SubmitObjectsRequest"`
		} `xml:"Message"`
	} `xml:"NotificationMessage"`
}

type DSUB_Interface interface {
	newEvent() error
}

func New_Transaction(i DSUB_Interface) error {
	return i.newEvent()
}

func (i *DSUBEvent) newEvent() error {
	var err error
	switch i.Action {
	case tukcnst.SELECT:
		log.Printf("Processing Select Subscriptions for Pathway %s", i.Pathway)
		err = i.selectSubscriptions()
	case tukcnst.CREATE:
		log.Println("Processing Create Subscriptions")
		err = i.createSubscriptions()
	case tukcnst.CANCEL:
		log.Println("Processing Cancel Subscriptions")
		err = i.cancelSubscriptions()
	default:
		log.Printf("Processing Broker Notify Message\n%s", i.EventMessage)
		i.processBrokerEventMessage()
	}
	return err
}

// creates a DSUBNotifyMessage from the EventMessage and populates a new TUKEvent with the DSUBNotifyMessage values
// It then checks for TukDB Subscriptions matching the brokerref and creates a TUKEvent for each subscription
// A DSUB ack response is always returned regardless of success
func (i *DSUBEvent) processBrokerEventMessage() {
	i.Response = []byte(tukcnst.GO_TEMPLATE_DSUB_ACK)
	if err := i.newDSUBNotifyMessage(); err == nil {
		if i.Event.BrokerRef == "" {
			log.Println("no subscription ref found in notification message")
			return
		}
		log.Printf("Found Subscription Reference %s Setting Event state from Notify Message", i.Event.BrokerRef)
		i.initTUKEvent()
		if i.Event.XdsPid == "" {
			log.Println("no pid found in notification message")
			return
		}
		log.Printf("Checking for TUK Event subscriptions with Broker Ref = %s", i.Event.BrokerRef)
		tukdbSubs := tukdbint.Subscriptions{Action: tukcnst.SELECT}
		tukdbSub := tukdbint.Subscription{BrokerRef: i.Event.BrokerRef}
		tukdbSubs.Subscriptions = append(tukdbSubs.Subscriptions, tukdbSub)
		if err = tukdbint.NewDBEvent(&tukdbSubs); err == nil {
			log.Printf("TUK Event Subscriptions Count : %v", tukdbSubs.Count)
			if tukdbSubs.Count > 0 {
				log.Printf("Obtaining NHS ID. Using %s", i.Event.XdsPid+":"+i.REG_OID)
				pdq := tukpdq.PDQQuery{
					Server_Mode: i.PDQ_SERVER_TYPE,
					REG_ID:      i.Event.XdsPid,
					Server_URL:  i.PDQ_SERVER_URL,
					REG_OID:     i.REG_OID,
				}
				if err = tukpdq.New_Transaction(&pdq); err != nil {
					log.Println(err.Error())
					return
				}
				if pdq.Count == 0 {
					log.Println("no patient returned for pid " + i.Event.XdsPid)
					return
				}

				for _, dbsub := range tukdbSubs.Subscriptions {
					if dbsub.Id > 0 {
						log.Printf("Creating event for %s %s Subsription for Broker Ref %s", dbsub.Pathway, dbsub.Expression, dbsub.BrokerRef)
						i.Event.Pathway = dbsub.Pathway
						i.Event.Topic = dbsub.Topic
						i.Event.NhsId = pdq.NHS_ID
						i.Event.EventType = tukcnst.XDW_TASKEVENTTYPE_ATTACHMENT
						tukevs := tukdbint.Events{Action: tukcnst.INSERT}
						tukevs.Events = append(tukevs.Events, i.Event)
						if err = tukdbint.NewDBEvent(&tukevs); err == nil {
							log.Printf("Created TUK DB Event for Pathway %s Expression %s Broker Ref %s", i.Event.Pathway, i.Event.Expression, i.Event.BrokerRef)
						}
					}
				}
			} else {
				log.Printf("No Subscriptions found with brokerref = %s. Sending Cancel request to Broker", i.Event.BrokerRef)
				i.cancelSubscriptions()
			}
		}
	} else {
		log.Println(err.Error())
	}

}

// InitDSUBEvent initialise the DSUBEvent struc with values parsed from the DSUBNotifyMessage
func (i *DSUBEvent) initTUKEvent() {
	i.Event.Creationtime = tukutil.Time_Now()
	i.Event.DocName = i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Name.LocalizedString.Value
	i.Event.BrokerRef = i.Notify.NotificationMessage.SubscriptionReference.Address.Text
	i.setRepositoryUniqueId()
	for _, c := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Classification {
		val := c.Name.LocalizedString.Value
		switch c.ClassificationScheme {
		case tukcnst.URN_CLASS_CODE:
			i.Event.ClassCode = val
			log.Printf("Set Class Code %s", val)
		case tukcnst.URN_CONF_CODE:
			i.Event.ConfCode = val
			log.Printf("Set Conf Code %s", val)
		case tukcnst.URN_FORMAT_CODE:
			i.Event.FormatCode = val
			log.Printf("Set Format Code %s", val)
		case tukcnst.URN_FACILITY_CODE:
			i.Event.FacilityCode = val
			log.Printf("Set Facility Code %s", val)
		case tukcnst.URN_PRACTICE_CODE:
			i.Event.PracticeCode = val
			log.Printf("Set Practice Code %s", val)
		case tukcnst.URN_TYPE_CODE:
			i.Event.Expression = val
			log.Printf("Set Type Code %s", val)
		case tukcnst.URN_AUTHOR:
			for _, slot := range c.Slot {
				switch slot.Name {
				case tukcnst.AUTHOR_PERSON:
					if i.Event.User == "" {
						user := strings.ReplaceAll(slot.ValueList.Value[0], "^", " ")
						user = strings.TrimSpace(user)
						i.Event.User = i.Event.User + user + " "
						log.Printf("Set User %s", i.Event.User)
					}
				case tukcnst.AUTHOR_INSTITUTION:
					if i.Event.Org == "" {
						i.Event.Org = strings.TrimSuffix(slot.ValueList.Value[0], "^^^")
						log.Printf("Set Org %s", i.Event.Org)
					}
				case tukcnst.AUTHOR_ROLE:
					if i.Event.Role == "" {
						i.Event.Role = strings.TrimSuffix(slot.ValueList.Value[0], "^^^")
						log.Printf("Set Role %s", i.Event.Role)
					}
				case tukcnst.AUTHOR_SPECIALITY:
					for _, value := range slot.ValueList.Value {
						if !strings.HasPrefix(i.Event.Speciality, value+". ") {
							i.Event.Speciality = i.Event.Speciality + strings.TrimSuffix(value, "^^^") + ". "
						}
					}
					log.Printf("Set Speciality %s", i.Event.Speciality)
					i.Event.Speciality = strings.TrimSpace(i.Event.Speciality)
				}
			}
		case tukcnst.URN_EVENT_LIST:
			i.Event.Comments = i.Event.Comments + val + ". "
		default:
			log.Printf("Unknown classication scheme %s. Skipping", c.ClassificationScheme)
		}
	}
	i.Event.Comments = strings.TrimSpace(i.Event.Comments)
	log.Printf("Added Event Codes to Notes %s", i.Event.Comments)
	i.setExternalIdentifiers()
	log.Println("Parsed DSUB Notify Message")
}

// NewDSUBNotifyMessage creates an IHE DSUB Notify message struc from the notfy element in the i.Message
func (i *DSUBEvent) newDSUBNotifyMessage() error {
	if i.EventMessage == "" {
		return errors.New("message is empty")
	}
	if err := xml.Unmarshal([]byte(i.EventMessage), &i.Notify); err != nil {
		return err
	}
	i.Event = tukdbint.Event{BrokerRef: i.Notify.NotificationMessage.SubscriptionReference.Address.Text}
	return nil
}
func (i *DSUBEvent) newDSUBCancelMessage() {
	tmplt, err := template.New(tukcnst.CANCEL).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_TEMPLATE_DSUB_CANCEL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var b bytes.Buffer
	err = tmplt.Execute(&b, i.BrokerRef)
	if err != nil {
		log.Println(err.Error())
		return
	}
	soapReq := tukhttp.SOAPRequest{
		URL:        i.BrokerURL,
		SOAPAction: tukcnst.SOAP_ACTION_UNSUBSCRIBE_REQUEST,
		Body:       b.Bytes(),
		Timeout:    2,
	}
	log.Printf("Sending Cancel Request to DSUB Broker %s", i.BrokerURL)
	err = tukhttp.NewRequest(&soapReq)
	if err != nil {
		log.Println(err.Error())
	}
}
func (i *DSUBEvent) setRepositoryUniqueId() {
	for _, slot := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Slot {
		if slot.Name == tukcnst.REPOSITORY_UID {
			i.Event.RepositoryUniqueId = slot.ValueList.Value[0]
			log.Printf("Set Repository UID %s", i.Event.RepositoryUniqueId)
			return
		}
	}
}
func (i *DSUBEvent) setExternalIdentifiers() {
	for exid := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.ExternalIdentifier {
		val := i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.ExternalIdentifier[exid].Value
		ids := i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.ExternalIdentifier[exid].IdentificationScheme
		switch ids {
		case tukcnst.URN_XDS_PID:
			i.Event.XdsPid = strings.Split(val, "^^^")[0]
			log.Printf("Set XDS PID %s", i.Event.XdsPid)
		case tukcnst.URN_XDS_DOCUID:
			i.Event.XdsDocEntryUid = val
			log.Printf("Set XDS Doc UID %s", i.Event.XdsDocEntryUid)
		}
	}
}

func (i *DSUBEvent) cancelSubscriptions() error {
	if i.Pathway == "" && i.RowID == 0 {
		log.Println("pathway or rowid not set. Sending Cancel subscription message to Broker")
		i.newDSUBCancelMessage()
		return nil
	}
	i.Subs = tukdbint.Subscriptions{Action: tukcnst.DELETE}
	delsub := tukdbint.Subscription{}
	if i.Pathway != "" {
		delsub.Pathway = i.Pathway
	} else {
		delsub.Id = i.RowID
	}
	i.Subs.Subscriptions = append(i.Subs.Subscriptions, delsub)
	return tukdbint.NewDBEvent(&i.Subs)
}

// (i *DSUBSubscribe) NewEvent() creates an IHE DSUB Subscribe request to an IHE DSUB broker. The Subscription reference is persisted in the subscriptions table.
func (i *DSUBEvent) createSubscriptions() error {
	type subreq struct {
		BrokerURL   string
		ConsumerURL string
		Topic       string
		Expression  string
	}
	i.Subs = tukdbint.Subscriptions{Action: i.Action}
	if len(i.Expressions) > 0 {
		for _, expression := range i.Expressions {
			log.Printf("Checking for existing Subscriptions for Expression %s", expression)
			expressionSubs := tukdbint.Subscriptions{Action: tukcnst.SELECT}
			expressionSub := tukdbint.Subscription{Expression: expression}
			expressionSubs.Subscriptions = append(expressionSubs.Subscriptions, expressionSub)
			if err := tukdbint.NewDBEvent(&expressionSubs); err != nil {
				log.Println(err.Error())
				return err
			}
			if expressionSubs.Count > 0 {
				log.Printf("Found %v Subscription(s) for Expression %s", expressionSubs.Count, expression)
				foundBrokerRef := false
				for _, s := range expressionSubs.Subscriptions {
					if !foundBrokerRef && s.BrokerRef != "" {
						log.Printf("Using DSUB Broker Reference %s for Pathway %s subscription to expression %s", expressionSubs.Subscriptions[1].BrokerRef, i.Pathway, expression)
						newSub := tukdbint.Subscription{}
						newSub.BrokerRef = s.BrokerRef
						log.Printf("Set Broker Ref =  %s", newSub.BrokerRef)
						newSub.Pathway = i.Pathway
						newSub.Expression = expression
						newSub.Topic = tukcnst.DSUB_TOPIC_TYPE_CODE
						newsubs := tukdbint.Subscriptions{Action: tukcnst.INSERT}
						newsubs.Subscriptions = append(newsubs.Subscriptions, newSub)
						log.Println("Registering Subscription with Event Service")
						if err := tukdbint.NewDBEvent(&newsubs); err != nil {
							log.Println(err.Error())
							return err
						}
						i.Subs.Subscriptions = append(i.Subs.Subscriptions, newSub)
						if i.Subs.LastInsertId < int64(newSub.Id) {
							i.Subs.LastInsertId = int64(newSub.Id)
						}
						i.Subs.Count = i.Subs.Count + 1
						foundBrokerRef = true
						log.Println("Registered Subscription with Event Service")
					}
				}
			} else {
				log.Printf("No Subscription for Expression %s found", expression)
				brokerSub := subreq{
					BrokerURL:   i.BrokerURL,
					ConsumerURL: i.ConsumerURL,
					Topic:       tukcnst.DSUB_TOPIC_TYPE_CODE,
					Expression:  expression,
				}
				tmplt, err := template.New(tukcnst.SUBSCRIBE).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_TEMPLATE_DSUB_SUBSCRIBE)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				var b bytes.Buffer
				err = tmplt.Execute(&b, brokerSub)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				soapReq := tukhttp.SOAPRequest{
					URL:        i.BrokerURL,
					SOAPAction: tukcnst.SOAP_ACTION_SUBSCRIBE_REQUEST,
					Body:       b.Bytes(),
				}
				log.Printf("Sending Subscribe Request to DSUB Broker %s", i.BrokerURL)
				err = tukhttp.NewRequest(&soapReq)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				subrsp := DSUBSubscribeResponse{}
				err = xml.Unmarshal(soapReq.Response, &subrsp)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				newSub := tukdbint.Subscription{}
				newSub.BrokerRef = subrsp.Body.SubscribeResponse.SubscriptionReference.Address
				log.Printf("Set Broker Ref =  %s", newSub.BrokerRef)
				if newSub.BrokerRef != "" {
					newSub.Pathway = i.Pathway
					newSub.Expression = brokerSub.Expression
					newSub.Topic = tukcnst.DSUB_TOPIC_TYPE_CODE
					newsubs := tukdbint.Subscriptions{Action: tukcnst.INSERT}
					newsubs.Subscriptions = append(newsubs.Subscriptions, newSub)
					log.Println("Registering Subscription with Event Service")
					if err := tukdbint.NewDBEvent(&newsubs); err != nil {
						log.Println(err.Error())
						return err
					}
					i.Subs.Subscriptions = append(i.Subs.Subscriptions, newSub)
					if i.Subs.LastInsertId < int64(newSub.Id) {
						i.Subs.LastInsertId = int64(newSub.Id)
					}
					i.Subs.Count = i.Subs.Count + 1
				} else {
					log.Println("No Broker Reference found. Subscription not registered with Event Service")
				}
			}
		}
	}
	i.Response, _ = json.MarshalIndent(i.Subs, "", "  ")
	return nil
}
func (i *DSUBEvent) selectSubscriptions() error {
	i.Subs = tukdbint.Subscriptions{Action: tukcnst.SELECT}
	sub := tukdbint.Subscription{}
	if i.Pathway != "" {
		sub.Pathway = i.Pathway
	}
	i.Subs.Subscriptions = append(i.Subs.Subscriptions, sub)
	if err := tukdbint.NewDBEvent(&i.Subs); err != nil {
		return err
	}
	i.Response, _ = json.MarshalIndent(i.Subs, "", "  ")
	return nil
}
