package tukxdw

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukdsub"
	"github.com/ipthomas/tukutil"
)

var DebugMode = true
var Regional_OID = os.Getenv(tukcnst.ENV_REG_OID)

type Interface interface {
	execute() error
}
type Transaction struct {
	Actor              string
	User               string
	Org                string
	Role               string
	Pathway            string
	Expression         string
	NHS_ID             string
	XDS_ID             string
	Task_ID            int
	XDWVersion         int
	Status             string
	AttachmentInfo     string
	Request            []byte
	Response           []byte
	ServiceURL         ServiceURL
	Dashboard          Dashboard
	WorkflowDefinition WorkflowDefinition
	XDSDocumentMeta    XDSDocumentMeta
	WorkflowDocument   WorkflowDocument
	XDWState           tukdbint.WorkflowStates
	Workflows          tukdbint.Workflows
	OpenWorkflows      tukdbint.Workflows
	ClosedWorkflows    tukdbint.Workflows
	MetWorkflows       tukdbint.Workflows
	OverdueWorkflows   tukdbint.Workflows
	EscalatedWorkflows tukdbint.Workflows
	XDWEvents          tukdbint.Events
	Subscriptions      tukdbint.Subscriptions
}
type ServiceURL struct {
	DSUB_Broker      string
	DSUB_Consumer    string
	Event_Consumer   string
	Event_Subscriber string
	Event_Notifier   string
	XDW_Consumer     string
	XDW_Creator      string
	HTML_Creator     string
	Service_Admin    string
}
type Dashboard struct {
	Total        int
	InProgress   int
	TargetMet    int
	TargetMissed int
	Escalated    int
	Complete     int
}
type XDSDocumentMeta struct {
	ID                    string `json:"id"`
	Repositoryuniqueid    string `json:"repositoryuniqueid"`
	Registryoid           string `json:"registryoid"`
	Languagecode          string `json:"languagecode"`
	Docname               string `json:"docname"`
	Docdesc               string `json:"docdesc"`
	DocID                 string `json:"docid"`
	Authorinstitution     string `json:"authorinstitution"`
	Authorperson          string `json:"authorperson"`
	Classcode             string `json:"classcode"`
	Classcodescheme       string `json:"classcodescheme"`
	Classcodevalue        string `json:"classcodevalue"`
	Typecode              string `json:"typecode"`
	Typecodescheme        string `json:"typecodescheme"`
	Typecodevalue         string `json:"typecodevalue"`
	Practicesettingcode   string `json:"practicesettingcode"`
	Practicesettingscheme string `json:"practicesettingscheme"`
	Practicesettingvalue  string `json:"practicesettingvalue"`
	Confcode              string `json:"confcode"`
	Confcodescheme        string `json:"confcodescheme"`
	Confcodevalue         string `json:"confcodevalue"`
	Facilitycode          string `json:"facilitycode"`
	Facilitycodescheme    string `json:"facilitycodescheme"`
	Facilitycodevalue     string `json:"facilitycodevalue"`
	Formatcode            string `json:"formatcode"`
	Formatcodescheme      string `json:"formatcodescheme"`
	Formatcodevalue       string `json:"formatcodevalue"`
	Mimetype              string `json:"mimetype"`
	Objecttype            string `json:"objecttype"`
}
type WorkflowDefinition struct {
	Ref                 string `json:"ref"`
	Name                string `json:"name"`
	Confidentialitycode string `json:"confidentialitycode"`
	StartByTime         string `json:"startbytime"`
	CompleteByTime      string `json:"completebytime"`
	ExpirationTime      string `json:"expirationtime"`
	CompletionBehavior  []struct {
		Completion struct {
			Condition string `json:"condition"`
		} `json:"completion"`
	} `json:"completionBehavior"`
	Tasks []struct {
		ID              string `json:"id"`
		Tasktype        string `json:"tasktype"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		ActualOwner     string `json:"actualowner"`
		ExpirationTime  string `json:"expirationtime,omitempty"`
		StartByTime     string `json:"startbytime,omitempty"`
		CompleteByTime  string `json:"completebytime"`
		IsSkipable      bool   `json:"isskipable,omitempty"`
		PotentialOwners []struct {
			OrganizationalEntity struct {
				User string `json:"user"`
			} `json:"organizationalEntity"`
		} `json:"potentialOwners,omitempty"`
		CompletionBehavior []struct {
			Completion struct {
				Condition string `json:"condition"`
			} `json:"completion"`
		} `json:"completionBehavior"`
		Input []struct {
			Name        string `json:"name"`
			Contenttype string `json:"contenttype"`
			AccessType  string `json:"accesstype"`
		} `json:"input,omitempty"`
		Output []struct {
			Name        string `json:"name"`
			Contenttype string `json:"contenttype"`
			AccessType  string `json:"accesstype"`
		} `json:"output,omitempty"`
	} `json:"tasks"`
}
type WorkflowDocument struct {
	XMLName                        xml.Name              `xml:"XDW.WorkflowDocument" json:"XMLName"`
	Hl7                            string                `xml:"hl7,attr" json:"Hl7"`
	WsHt                           string                `xml:"ws-ht,attr" json:"WsHt"`
	Xdw                            string                `xml:"xdw,attr" json:"Xdw"`
	Xsi                            string                `xml:"xsi,attr" json:"Xsi"`
	SchemaLocation                 string                `xml:"schemaLocation,attr" json:"SchemaLocation"`
	ID                             ID                    `xml:"id" json:"ID"`
	EffectiveTime                  EffectiveTime         `xml:"effectiveTime" json:"EffectiveTime"`
	ConfidentialityCode            ConfidentialityCode   `xml:"confidentialityCode" json:"ConfidentialityCode"`
	Patient                        ID                    `xml:"patient" json:"Patient"`
	Author                         Author                `xml:"author" json:"Author"`
	WorkflowInstanceId             string                `xml:"workflowInstanceId" json:"WorkflowInstanceId"`
	WorkflowDocumentSequenceNumber string                `xml:"workflowDocumentSequenceNumber" json:"WorkflowDocumentSequenceNumber"`
	WorkflowStatus                 string                `xml:"workflowStatus" json:"WorkflowStatus"`
	WorkflowStatusHistory          WorkflowStatusHistory `xml:"workflowStatusHistory" json:"WorkflowStatusHistory"`
	WorkflowDefinitionReference    string                `xml:"workflowDefinitionReference" json:"WorkflowDefinitionReference"`
	TaskList                       TaskList              `xml:"TaskList" json:"TaskList"`
}
type ConfidentialityCode struct {
	Code string `xml:"code,attr" json:"Code"`
}
type EffectiveTime struct {
	Value string `xml:"value,attr" json:"Value"`
}
type Author struct {
	AssignedAuthor AssignedAuthor `xml:"assignedAuthor" json:"AssignedAuthor"`
}
type AssignedAuthor struct {
	ID             ID             `xml:"id" json:"ID"`
	AssignedPerson AssignedPerson `xml:"assignedPerson" json:"AssignedPerson"`
}
type ID struct {
	Root                   string `xml:"root,attr" json:"Root"`
	Extension              string `xml:"extension,attr" json:"Extension"`
	AssigningAuthorityName string `xml:"assigningAuthorityName,attr" json:"AssigningAuthorityName"`
}
type AssignedPerson struct {
	Name Name `xml:"name" json:"Name"`
}
type Name struct {
	Family string `xml:"family" json:"Family"`
	Prefix string `xml:"prefix" json:"Prefix"`
}
type WorkflowStatusHistory struct {
	DocumentEvent []DocumentEvent `xml:"documentEvent" json:"DocumentEvent"`
}
type TaskList struct {
	XDWTask []XDWTask `xml:"XDWTask" json:"Task"`
}
type XDWTask struct {
	TaskData         TaskData         `xml:"taskData" json:"TaskData"`
	TaskEventHistory TaskEventHistory `xml:"taskEventHistory" json:"TaskEventHistory"`
}
type TaskData struct {
	TaskDetails TaskDetails `xml:"taskDetails" json:"TaskDetails"`
	Description string      `xml:"description" json:"Description"`
	Input       []Input     `xml:"input" json:"Input"`
	Output      []Output    `xml:"output" json:"Output"`
}
type TaskDetails struct {
	ID                    string `xml:"id" json:"ID"`
	TaskType              string `xml:"taskType" json:"TaskType"`
	Name                  string `xml:"name" json:"Name"`
	Status                string `xml:"status" json:"Status"`
	ActualOwner           string `xml:"actualOwner" json:"ActualOwner"`
	CreatedTime           string `xml:"createdTime" json:"CreatedTime"`
	CreatedBy             string `xml:"createdBy" json:"CreatedBy"`
	ActivationTime        string `xml:"activationTime" json:"ActivationTime"`
	LastModifiedTime      string `xml:"lastModifiedTime" json:"LastModifiedTime"`
	RenderingMethodExists string `xml:"renderingMethodExists" json:"RenderingMethodExists"`
}
type TaskEventHistory struct {
	TaskEvent []TaskEvent `xml:"taskEvent" json:"TaskEvent"`
}
type AttachmentInfo struct {
	Identifier      string `xml:"identifier" json:"Identifier"`
	Name            string `xml:"name" json:"Name"`
	AccessType      string `xml:"accessType" json:"AccessType"`
	ContentType     string `xml:"contentType" json:"ContentType"`
	ContentCategory string `xml:"contentCategory" json:"ContentCategory"`
	AttachedTime    string `xml:"attachedTime" json:"AttachedTime"`
	AttachedBy      string `xml:"attachedBy" json:"AttachedBy"`
	HomeCommunityId string `xml:"homeCommunityId" json:"HomeCommunityId"`
}
type Part struct {
	Name           string         `xml:"name,attr" json:"Name"`
	AttachmentInfo AttachmentInfo `xml:"attachmentInfo" json:"AttachmentInfo"`
}
type Output struct {
	Part Part `xml:"part" json:"Part"`
}
type Input struct {
	Part Part `xml:"part" json:"Part"`
}
type DocumentEvent struct {
	EventTime           string `xml:"eventTime" json:"EventTime"`
	EventType           string `xml:"eventType" json:"EventType"`
	TaskEventIdentifier string `xml:"taskEventIdentifier" json:"TaskEventIdentifier"`
	Author              string `xml:"author" json:"Author"`
	PreviousStatus      string `xml:"previousStatus" json:"PreviousStatus"`
	ActualStatus        string `xml:"actualStatus" json:"ActualStatus"`
}
type TaskEvent struct {
	ID         string `xml:"id" json:"ID"`
	EventTime  string `xml:"eventTime" json:"EventTime"`
	Identifier string `xml:"identifier" json:"Identifier"`
	EventType  string `xml:"eventType" json:"EventType"`
	Status     string `xml:"status" json:"Status"`
}

// sort interface for Document Events

func (e WorkflowStatusHistory) Len() int {
	return len(e.DocumentEvent)
}
func (e WorkflowStatusHistory) Less(i, j int) bool {
	return e.DocumentEvent[i].EventTime > e.DocumentEvent[j].EventTime
}
func (e WorkflowStatusHistory) Swap(i, j int) {
	e.DocumentEvent[i], e.DocumentEvent[j] = e.DocumentEvent[j], e.DocumentEvent[i]
}

// sort interface for Task Events

func (e TaskEventHistory) Len() int {
	return len(e.TaskEvent)
}
func (e TaskEventHistory) Less(i, j int) bool {
	return e.TaskEvent[i].EventTime > e.TaskEvent[j].EventTime
}
func (e TaskEventHistory) Swap(i, j int) {
	e.TaskEvent[i], e.TaskEvent[j] = e.TaskEvent[j], e.TaskEvent[i]
}

func Execute(i Interface) error {
	return i.execute()
}

// IHE XDW Actors

func (i *Transaction) execute() error {
	switch i.Actor {
	case tukcnst.XDW_ADMIN_REGISTER_DEFINITION:
		return i.RegisterWorkflowDefinition(false)
	case tukcnst.XDW_ADMIN_REGISTER_XDS_META:
		return i.RegisterWorkflowDefinition(true)
	case tukcnst.XDW_ACTOR_CONTENT_CREATOR:
		return i.ContentCreator()
	case tukcnst.XDW_ACTOR_CONTENT_CONSUMER:
		return i.ContentConsumer()
	case tukcnst.XDW_ACTOR_CONTENT_UPDATER:
		return i.ContentUpdater()
	}
	return nil
}

// IHE XDW Content Updater

func (i *Transaction) ContentUpdater() error {
	i.Workflows = tukdbint.Workflows{Action: tukcnst.SELECT}
	wf := tukdbint.Workflow{Pathway: i.Pathway, NHSId: i.NHS_ID, Version: i.XDWVersion}
	i.Workflows.Workflows = append(i.Workflows.Workflows, wf)
	if err := tukdbint.NewDBEvent(&i.Workflows); err != nil {
		log.Println(err.Error())
		return err
	}
	i.XDWEvents = tukdbint.Events{Action: tukcnst.SELECT}
	ev := tukdbint.Event{Pathway: i.Pathway, NhsId: i.NHS_ID, Version: i.XDWVersion, TaskId: -1}
	i.XDWEvents.Events = append(i.XDWEvents.Events, ev)
	if err := tukdbint.NewDBEvent(&i.XDWEvents); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Updating state of %v Workflows with %v Events", i.Workflows.Count, i.XDWEvents.Count)

	for _, wf := range i.Workflows.Workflows {
		if wf.Id == 0 {
			continue
		}
		log.Printf("Updating %s Workflow Version %v for NHS ID %s", wf.Pathway, wf.Version, wf.NHSId)

		i.XDWVersion = wf.Version
		if err := json.Unmarshal([]byte(wf.XDW_Def), &i.WorkflowDefinition); err != nil {
			log.Println(err.Error())
			return err
		}
		if err := json.Unmarshal([]byte(wf.XDW_Doc), &i.WorkflowDocument); err != nil {
			log.Println(err.Error())
			return err
		}
		log.Printf("Processing %v Events", i.XDWEvents.Count)
		newEvents := tukdbint.Events{}
		cnt := 0
		for _, ev := range i.XDWEvents.Events {
			if ev.Id == 0 {
				continue
			}
			cnt = cnt + 1
			log.Printf("Processing Event %v ID %v Obtaining Workflow Task", cnt, ev.Id)
			newevent := true
			for _, task := range i.WorkflowDocument.TaskList.XDWTask {
				if task.TaskData.TaskDetails.ID == tukutil.GetStringFromInt(ev.TaskId) {
					log.Printf("Searching Task %s eevents for matching event ID %v", task.TaskData.TaskDetails.ID, ev.Id)
					for _, taskevent := range task.TaskEventHistory.TaskEvent {
						if taskevent.ID == tukutil.GetStringFromInt(int(ev.Id)) {
							log.Printf("Task %s Event %v is registered. Skipping Event", task.TaskData.TaskDetails.ID, ev.Id)
							newevent = false
						}
					}
				}
			}
			if newevent {
				newEvents.Events = append(newEvents.Events, ev)
				log.Printf("Event %v is not registered. Including Event in Workflow Task events updates", ev.Id)
			}
		}
		if len(newEvents.Events) > 0 {
			log.Printf("Updating Workflow with %v new events", len(newEvents.Events))
			sort.Sort(eventsList(i.XDWEvents.Events))
			i.XDWEvents.Events = newEvents.Events
			i.XDWEvents.Count = len(newEvents.Events)
			if err := i.UpdateWorkflowDocumentTasks(); err != nil {
				log.Println(err.Error())
			}
		}
	}
	return nil
}
func (i *Transaction) UpdateWorkflowDocumentTasks() error {
	log.Printf("Updating %s Workflow Tasks with %v Events", i.WorkflowDocument.WorkflowDefinitionReference, len(i.XDWEvents.Events))
	for _, ev := range i.XDWEvents.Events {
		for k, wfdoctask := range i.WorkflowDocument.TaskList.XDWTask {
			log.Println("Checking Workflow Document Task " + wfdoctask.TaskData.TaskDetails.Name + " for matching Events")
			for inp, input := range wfdoctask.TaskData.Input {
				if ev.Expression == input.Part.Name {
					log.Println("Matched workflow document task " + wfdoctask.TaskData.TaskDetails.ID + " Input Part : " + input.Part.Name + " with Event Expression : " + ev.Expression + " Status : " + wfdoctask.TaskData.TaskDetails.Status)
					if !i.isInputRegistered(ev) {
						log.Printf("Updating XDW with Event ID %v for Task ID %s", ev.Id, wfdoctask.TaskData.TaskDetails.ID)
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Input[inp].Part.AttachmentInfo.AttachedTime = ev.Creationtime
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Input[inp].Part.AttachmentInfo.AttachedBy = ev.User + " " + ev.Org + " " + ev.Role
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Input[inp].Part.AttachmentInfo.HomeCommunityId = Regional_OID
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.LastModifiedTime = ev.Creationtime
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.Status = tukcnst.IN_PROGRESS
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActualOwner = ev.User + " " + ev.Org + " " + ev.Role
						if i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActivationTime == "" {
							i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActivationTime = ev.Creationtime
							log.Printf("Set Task %s Activation Time %s", wfdoctask.TaskData.TaskDetails.ID, i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActivationTime)
						}
						if wfdoctask.TaskData.Input[inp].Part.AttachmentInfo.AccessType == tukcnst.XDS_REGISTERED {
							i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Input[inp].Part.AttachmentInfo.Identifier = ev.XdsDocEntryUid
						} else {
							i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Input[inp].Part.AttachmentInfo.Identifier = tukutil.GetStringFromInt(int(ev.Id))
						}
						i.newTaskEvent(ev)
						wfseqnum, _ := strconv.ParseInt(i.WorkflowDocument.WorkflowDocumentSequenceNumber, 0, 0)
						wfseqnum = wfseqnum + 1
						i.WorkflowDocument.WorkflowDocumentSequenceNumber = strconv.Itoa(int(wfseqnum))
						i.newDocEvent(ev)
					}
				}
			}
			for oup, output := range i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Output {
				if ev.Expression == output.Part.Name {
					log.Println("Matched workflow document task " + wfdoctask.TaskData.TaskDetails.ID + " Output Part : " + output.Part.Name + " with Event Expression : " + ev.Expression + " Status : " + wfdoctask.TaskData.TaskDetails.Status)
					if !i.isOutputRegistered(ev) {
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.LastModifiedTime = ev.Creationtime
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Output[oup].Part.AttachmentInfo.AttachedTime = ev.Creationtime
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Output[oup].Part.AttachmentInfo.AttachedBy = ev.User + " " + ev.Org + " " + ev.Role
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActualOwner = ev.User + " " + ev.Org + " " + ev.Role
						i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.Status = tukcnst.IN_PROGRESS
						if i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActivationTime == "" {
							i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.ActivationTime = ev.Creationtime
						}
						if strings.HasSuffix(wfdoctask.TaskData.Output[oup].Part.AttachmentInfo.AccessType, tukcnst.XDS_REGISTERED) {
							i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Output[oup].Part.AttachmentInfo.Identifier = ev.XdsDocEntryUid
						} else {
							i.WorkflowDocument.TaskList.XDWTask[k].TaskData.Output[oup].Part.AttachmentInfo.Identifier = tukutil.GetStringFromInt(int(ev.Id))
						}
						i.newTaskEvent(ev)
						wfseqnum, _ := strconv.ParseInt(i.WorkflowDocument.WorkflowDocumentSequenceNumber, 0, 0)
						wfseqnum = wfseqnum + 1
						i.WorkflowDocument.WorkflowDocumentSequenceNumber = strconv.Itoa(int(wfseqnum))
						i.newDocEvent(ev)
					}
				}
			}
		}
	}

	for k, task := range i.WorkflowDocument.TaskList.XDWTask {
		i.Task_ID = tukutil.GetIntFromString(task.TaskData.TaskDetails.ID) - 1
		if i.IsTaskCompleteBehaviorMet() {
			i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.Status = tukcnst.COMPLETE
		}
	}

	if complete, completiontask := i.IsWorkflowCompleteBehaviorMet(); complete {
		i.WorkflowDocument.WorkflowStatus = tukcnst.CLOSED
		i.newEventID(tukcnst.XDW_TASKEVENTTYPE_WORKFLOW_COMPLETED, i.XDWVersion)
		docevent := DocumentEvent{}
		docevent.Author = i.User + " " + i.Org + " " + i.Role
		docevent.TaskEventIdentifier = completiontask
		docevent.EventTime = tukutil.Time_Now()
		docevent.EventType = tukcnst.XDW_DOCEVENTTYPE_COMPLETED_WORKFLOW
		docevent.PreviousStatus = i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent[len(i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent)-1].ActualStatus
		docevent.ActualStatus = tukcnst.COMPLETE
		i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent = append(i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent, docevent)
		for k := range i.WorkflowDocument.TaskList.XDWTask {
			i.WorkflowDocument.TaskList.XDWTask[k].TaskData.TaskDetails.Status = tukcnst.COMPLETE
		}
		log.Println("Closed Workflow. Total Workflow Document Events " + strconv.Itoa(len(i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent)))
	}
	return i.updateWorkflow()
}
func (i *Transaction) isInputRegistered(ev tukdbint.Event) bool {
	log.Printf("Checking if Input Event for Task %s is registered", i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskData.Description)
	for _, input := range i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskData.Input {
		if ev.Expression == input.Part.Name {
			if input.Part.AttachmentInfo.AccessType == tukcnst.XDS_REGISTERED {
				if input.Part.AttachmentInfo.Identifier == ev.XdsDocEntryUid {
					log.Println("Event is registered. Skipping Event ")
					return true
				}
			} else {
				if input.Part.AttachmentInfo.Identifier == tukutil.GetStringFromInt(int(ev.Id)) {
					log.Println("Event is registered. Skipping Event ")
					return true
				}
			}
		}
	}
	return false
}
func (i *Transaction) isOutputRegistered(ev tukdbint.Event) bool {
	log.Printf("Checking if Ouput Event for Task %s is registered", i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskData.Description)
	for _, output := range i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskData.Output {
		if ev.Expression == output.Part.Name {
			if output.Part.AttachmentInfo.AccessType == tukcnst.XDS_REGISTERED {
				if output.Part.AttachmentInfo.Identifier == ev.XdsDocEntryUid {
					log.Println("Event is registered. Skipping Event ")
					return true
				}
			} else {
				if output.Part.AttachmentInfo.Identifier == tukutil.GetStringFromInt(int(ev.Id)) {
					log.Println("Event is registered. Skipping Event ")
					return true
				}
			}
		}
	}
	return false
}
func (i *Transaction) newTaskEvent(ev tukdbint.Event) {
	nte := TaskEvent{
		ID:         tukutil.GetStringFromInt(int(ev.Id)),
		Identifier: tukutil.GetStringFromInt(ev.TaskId),
		EventType:  i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskData.TaskDetails.TaskType,
		Status:     tukcnst.COMPLETE,
	}
	i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskEventHistory.TaskEvent = append(i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskEventHistory.TaskEvent, nte)
}
func (i *Transaction) newDocEvent(ev tukdbint.Event) {
	docevent := DocumentEvent{}
	docevent.Author = ev.User + " " + ev.Org + " " + ev.Role
	docevent.TaskEventIdentifier = tukutil.GetStringFromInt(ev.TaskId)
	docevent.EventTime = ev.Creationtime
	docevent.EventType = i.WorkflowDocument.TaskList.XDWTask[ev.TaskId-1].TaskData.TaskDetails.TaskType
	docevent.PreviousStatus = i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent[len(i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent)-1].ActualStatus
	docevent.ActualStatus = tukcnst.IN_PROGRESS
	i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent = append(i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent, docevent)
}
func (i *Transaction) updateWorkflow() error {
	var err error
	wfs := tukdbint.Workflows{Action: tukcnst.UPDATE}
	wf := tukdbint.Workflow{
		Pathway: i.Pathway,
		NHSId:   i.NHS_ID,
		XDW_Key: strings.ToUpper(i.Pathway) + i.NHS_ID,
		XDW_UID: i.WorkflowDocument.ID.Extension,
		Version: i.XDWVersion,
		Status:  i.WorkflowDocument.WorkflowStatus,
	}

	xdwDocBytes, _ := json.MarshalIndent(i.WorkflowDocument, "", "  ")
	wf.XDW_Doc = string(xdwDocBytes)
	wfs.Workflows = append(wfs.Workflows, wf)
	if err = tukdbint.NewDBEvent(&wfs); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("Upddated Workflow State for Pathway %s NHS ID %s Version %v Status %s", i.Pathway, i.NHS_ID, i.XDWVersion, i.WorkflowDocument.WorkflowStatus)
	}
	return err
}

// IHE XDW Content Creator
func (i *Transaction) ContentCreator() error {
	log.Printf("Creating New Workflow for Pathway %s NHS ID %s", i.Pathway, i.NHS_ID)
	var err error
	if err = i.loadWorkflowDefinition(); err == nil {
		if err = i.deprecateWorkflow(); err == nil {
			i.createWorkflow()
			err = i.persistWorkflow()
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func (i *Transaction) loadWorkflowDefinition() error {
	var err error
	xdwsdef := tukdbint.XDWS{Action: tukcnst.SELECT}
	xdwdef := tukdbint.XDW{Name: i.Pathway}
	xdwsdef.XDW = append(xdwsdef.XDW, xdwdef)
	if err = tukdbint.NewDBEvent(&xdwsdef); err == nil && xdwsdef.Count == 1 {
		if err = json.Unmarshal([]byte(xdwsdef.XDW[1].XDW), &i.WorkflowDefinition); err == nil {
			log.Printf("Loaded XDW definition for Pathway %s", i.Pathway)
		}
	} else {
		err = errors.New("no xdw definition registered for workflow" + i.Pathway + " you need to register the " + i.Pathway + " workflow definition before you can create a new workflow")
		log.Println(err.Error())
	}
	return err
}
func (i *Transaction) deprecateWorkflow() error {
	var err error
	wfs := tukdbint.Workflows{Action: tukcnst.DEPRECATE}
	wf := tukdbint.Workflow{XDW_Key: i.Pathway + i.NHS_ID}
	wfs.Workflows = append(wfs.Workflows, wf)
	log.Printf("Deprecating any current %s Workflow for NHS ID %s", i.Pathway, i.NHS_ID)
	if err = tukdbint.NewDBEvent(&wfs); err == nil {
		evs := tukdbint.Events{Action: tukcnst.DEPRECATE}
		ev := tukdbint.Event{Pathway: i.Pathway, NhsId: i.NHS_ID}
		evs.Events = append(evs.Events, ev)
		log.Printf("Deprecating any current %s Workflow events for NHS ID %s", i.Pathway, i.NHS_ID)
		if err = tukdbint.NewDBEvent(&evs); err != nil {
			log.Println(err.Error())
		}
	}
	return err
}
func (i *Transaction) createWorkflow() {
	i.Expression = tukcnst.XDW_DOCEVENTTYPE_CREATE_WORKFLOW
	var authoid = getLocalId(i.Org)
	var patoid = tukcnst.NHS_OID_DEFAULT
	var wfid = tukutil.Newid()

	var effectiveTime = tukutil.Time_Now()
	i.WorkflowDocument.Xdw = tukcnst.XDWNameSpace
	i.WorkflowDocument.Hl7 = tukcnst.HL7NameSpace
	i.WorkflowDocument.WsHt = tukcnst.WHTNameSpace
	i.WorkflowDocument.Xsi = tukcnst.XMLNS_XSI
	i.WorkflowDocument.XMLName.Local = tukcnst.XDWNameLocal
	i.WorkflowDocument.SchemaLocation = tukcnst.WorkflowDocumentSchemaLocation
	i.WorkflowDocument.ID.Root = strings.ReplaceAll(tukcnst.WorkflowInstanceId, "^", "")
	i.WorkflowDocument.ID.Extension = wfid
	i.WorkflowDocument.ID.AssigningAuthorityName = strings.ToUpper(i.Org)
	i.WorkflowDocument.EffectiveTime.Value = effectiveTime
	i.WorkflowDocument.ConfidentialityCode.Code = i.WorkflowDefinition.Confidentialitycode
	i.WorkflowDocument.Patient.Root = patoid
	i.WorkflowDocument.Patient.Extension = i.NHS_ID
	i.WorkflowDocument.Patient.AssigningAuthorityName = "NHS"
	i.WorkflowDocument.Author.AssignedAuthor.ID.Root = authoid
	i.WorkflowDocument.Author.AssignedAuthor.ID.Extension = strings.ToUpper(i.Org)
	i.WorkflowDocument.Author.AssignedAuthor.ID.AssigningAuthorityName = authoid
	i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Family = i.User
	i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Prefix = i.Role
	i.WorkflowDocument.WorkflowInstanceId = wfid + tukcnst.WorkflowInstanceId
	i.WorkflowDocument.WorkflowDocumentSequenceNumber = "1"
	i.WorkflowDocument.WorkflowStatus = tukcnst.OPEN
	i.WorkflowDocument.WorkflowDefinitionReference = strings.ToUpper(i.Pathway)
	var tevidstr = tukutil.GetStringFromInt(int(i.newEventID(tukcnst.XDW_DOCEVENTTYPE_CREATE_WORKFLOW, i.XDWVersion)))

	docevent := DocumentEvent{}
	docevent.Author = i.User + " " + i.Org + " " + i.Role
	docevent.TaskEventIdentifier = "0"
	docevent.EventTime = effectiveTime
	docevent.EventType = tukcnst.XDW_DOCEVENTTYPE_CREATE_WORKFLOW
	docevent.ActualStatus = tukcnst.OPEN
	i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent = append(i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent, docevent)

	for _, t := range i.WorkflowDefinition.Tasks {
		i.Expression = t.Name
		i.Task_ID = tukutil.GetIntFromString(t.ID)
		tevidstr = tukutil.GetStringFromInt(int(i.newEventID(tukcnst.XDW_TASKEVENTTYPE_CREATE_TASK, i.XDWVersion)))
		log.Printf("Creating Workflow Task ID - %v Name - %s", t.ID, t.Name)
		task := XDWTask{}
		task.TaskData.TaskDetails.ID = t.ID
		task.TaskData.TaskDetails.TaskType = t.Tasktype
		task.TaskData.TaskDetails.Name = t.Name
		task.TaskData.TaskDetails.ActualOwner = t.ActualOwner
		task.TaskData.TaskDetails.CreatedBy = i.Role + " " + i.User
		task.TaskData.TaskDetails.CreatedTime = effectiveTime
		task.TaskData.TaskDetails.RenderingMethodExists = "false"
		task.TaskData.TaskDetails.LastModifiedTime = effectiveTime
		task.TaskData.Description = t.Description
		task.TaskData.TaskDetails.Status = tukcnst.CREATED
		for _, inp := range t.Input {
			docinput := Input{}
			docinput.Part.Name = inp.Name
			docinput.Part.AttachmentInfo.Name = inp.Name
			docinput.Part.AttachmentInfo.AccessType = inp.AccessType
			docinput.Part.AttachmentInfo.ContentType = inp.Contenttype
			docinput.Part.AttachmentInfo.ContentCategory = tukcnst.MEDIA_TYPES
			task.TaskData.Input = append(task.TaskData.Input, docinput)
			log.Printf("Created Input Part - %s", inp.Name)
		}
		for _, outp := range t.Output {
			docoutput := Output{}
			docoutput.Part.Name = outp.Name
			docoutput.Part.AttachmentInfo.Name = outp.Name
			docoutput.Part.AttachmentInfo.AccessType = outp.AccessType
			docoutput.Part.AttachmentInfo.ContentType = outp.Contenttype
			docoutput.Part.AttachmentInfo.ContentCategory = tukcnst.MEDIA_TYPES
			task.TaskData.Output = append(task.TaskData.Output, docoutput)
			log.Printf("Created Output Part - %s", outp.Name)
		}
		tev := TaskEvent{}
		tev.EventTime = effectiveTime
		tev.ID = tevidstr
		tev.Identifier = t.ID
		tev.EventType = tukcnst.XDW_TASKEVENTTYPE_CREATED
		tev.Status = tukcnst.XDW_TASKEVENTTYPE_COMPLETE
		task.TaskEventHistory.TaskEvent = append(task.TaskEventHistory.TaskEvent, tev)
		i.WorkflowDocument.TaskList.XDWTask = append(i.WorkflowDocument.TaskList.XDWTask, task)
		log.Printf("Set Workflow Task Event %s %s status to %s", t.ID, tev.EventType, tev.Status)
	}
	i.Response, _ = json.MarshalIndent(i.WorkflowDocument, "", "  ")
	i.XDWVersion = 0
	log.Printf("%s Created new %s Workflow for Patient %s", i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Family, i.WorkflowDocument.WorkflowDefinitionReference, i.NHS_ID)
}
func (i *Transaction) persistWorkflow() error {
	var err error
	wfs := tukdbint.Workflows{Action: tukcnst.INSERT}
	wf := tukdbint.Workflow{
		Pathway: i.Pathway,
		NHSId:   i.NHS_ID,
		XDW_Key: strings.ToUpper(i.Pathway) + i.NHS_ID,
		XDW_UID: i.WorkflowDocument.ID.Extension,
		Version: i.XDWVersion,
		Status:  i.WorkflowDocument.WorkflowStatus,
	}
	xdwDocBytes, _ := json.MarshalIndent(i.WorkflowDocument, "", "  ")
	xdwDefBytes, _ := json.Marshal(i.WorkflowDefinition)
	wf.XDW_Doc = string(xdwDocBytes)
	wf.XDW_Def = string(xdwDefBytes)
	wfs.Workflows = append(wfs.Workflows, wf)
	if err = tukdbint.NewDBEvent(&wfs); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("Persisted Workflow Version %v for Pathway %s NHS ID %s", i.XDWVersion, i.Pathway, i.NHS_ID)
	}
	i.Workflows.LastInsertId = wfs.LastInsertId
	return err
}

// IHE XDW Content Consumer
func (i *Transaction) ContentConsumer() error {
	i.Workflows = tukdbint.Workflows{Action: tukcnst.SELECT}
	wf := tukdbint.Workflow{Pathway: i.Pathway, NHSId: i.NHS_ID, Version: i.XDWVersion, Status: i.Status}
	i.Workflows.Workflows = append(i.Workflows.Workflows, wf)
	if err := tukdbint.NewDBEvent(&i.Workflows); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Selected %v Workflows", i.Workflows.Count)
	if err := i.SetXDWStates(); err != nil {
		log.Println(err.Error())
		return err
	}
	i.XDWEvents = tukdbint.Events{Action: tukcnst.SELECT}
	ev := tukdbint.Event{Pathway: i.Pathway, NhsId: i.NHS_ID, Version: i.XDWVersion, TaskId: -1}
	i.XDWEvents.Events = append(i.XDWEvents.Events, ev)
	if err := tukdbint.NewDBEvent(&i.XDWEvents); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
func (i *Transaction) SetXDWStates() error {
	log.Println("Setting XDW States")
	var err error
	i.Dashboard.Total = i.Workflows.Count

	for _, wf := range i.Workflows.Workflows {
		if len(wf.XDW_Doc) > 0 {
			if err = json.Unmarshal([]byte(wf.XDW_Doc), &i.WorkflowDocument); err != nil {
				log.Println(err.Error())
				return err
			}
			if err = json.Unmarshal([]byte(wf.XDW_Def), &i.WorkflowDefinition); err != nil {
				log.Println(err.Error())
				return err
			}
			log.Printf("Setting %s Workflow state for Patient %s", i.WorkflowDocument.WorkflowDefinitionReference, i.WorkflowDocument.Patient.Extension)
			state := tukdbint.Workflowstate{}
			state.Created = wf.Created
			state.Status = wf.Status
			state.Published = wf.Published
			state.WorkflowId = wf.Id
			state.Pathway = wf.Pathway
			state.NHSId = wf.NHSId
			state.Version = wf.Version
			state.CreatedBy = i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Family + " " + i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Prefix
			state.CompleteBy = "Non Specified"
			state.LastUpdate = i.GetLatestWorkflowEventTime().String()
			state.Owner = ""
			state.Overdue = "FALSE"
			state.Escalated = "FALSE"
			state.TargetMet = "TRUE"
			state.InProgress = "TRUE"
			state.Duration = i.GetWorkflowDuration()
			if state.Status == tukcnst.TUK_STATUS_CLOSED {
				state.TimeRemaining = "0"
				i.ClosedWorkflows.Workflows = append(i.ClosedWorkflows.Workflows, wf)
				i.ClosedWorkflows.Count = i.ClosedWorkflows.Count + 1
			} else {
				state.TimeRemaining = i.GetWorkflowTimeRemaining()
				i.OpenWorkflows.Workflows = append(i.OpenWorkflows.Workflows, wf)
				i.OpenWorkflows.Count = i.OpenWorkflows.Count + 1
			}
			workflowStartTime := tukutil.GetTimeFromString(state.Created)
			if i.IsWorkflowOverdue() {
				i.Dashboard.TargetMissed = i.Dashboard.TargetMissed + 1
				state.Overdue = "TRUE"
				state.TargetMet = "FALSE"
				i.OverdueWorkflows.Workflows = append(i.OverdueWorkflows.Workflows, wf)
				i.OverdueWorkflows.Count = i.OverdueWorkflows.Count + 1
			} else {
				if i.WorkflowDocument.WorkflowStatus == tukcnst.CLOSED {
					i.Dashboard.TargetMet = i.Dashboard.TargetMet + 1
					i.MetWorkflows.Workflows = append(i.MetWorkflows.Workflows, wf)
					i.MetWorkflows.Count = i.MetWorkflows.Count + 1
				}
			}
			if i.WorkflowDefinition.CompleteByTime == "" {
				state.CompleteBy = "Non Specified"
			} else {
				period := strings.Split(i.WorkflowDefinition.CompleteByTime, "(")[0]
				periodDuration := tukutil.GetIntFromString(strings.Split(strings.Split(i.WorkflowDefinition.CompleteByTime, "(")[1], ")")[0])
				switch period {
				case "month":
					state.CompleteBy = strings.Split(tukutil.GetFutureDate(workflowStartTime, 0, periodDuration, 0, 0, 0).String(), " +")[0]
				case "day":
					state.CompleteBy = strings.Split(tukutil.GetFutureDate(workflowStartTime, 0, 0, periodDuration, 0, 0).String(), " +")[0]
				case "hour":
					state.CompleteBy = strings.Split(tukutil.GetFutureDate(workflowStartTime, 0, 0, 0, periodDuration, 0).String(), " +")[0]
				case "min":
					state.CompleteBy = strings.Split(tukutil.GetFutureDate(workflowStartTime, 0, 0, 0, 0, periodDuration).String(), " +")[0]
				}
			}

			if i.WorkflowDocument.WorkflowStatus == tukcnst.OPEN {
				log.Printf("Workflow %s is OPEN", wf.XDW_Key)
				i.Dashboard.InProgress = i.Dashboard.InProgress + 1
				if i.IsWorkflowEscalated() {
					log.Printf("Workflow %s is ESCALATED", wf.XDW_Key)
					i.Dashboard.Escalated = i.Dashboard.Escalated + 1
					state.Escalated = "TRUE"
					i.EscalatedWorkflows.Workflows = append(i.EscalatedWorkflows.Workflows, wf)
					i.EscalatedWorkflows.Count = i.EscalatedWorkflows.Count + 1
				}
			} else {
				log.Printf("Workflow %s is CLOSED", wf.XDW_Key)
				i.Dashboard.Complete = i.Dashboard.Complete + 1
				state.InProgress = "FALSE"
			}
			i.XDWState.Workflowstate = append(i.XDWState.Workflowstate, state)
		}
	}

	return err
}
func (i *Transaction) GetRegisteredWorkflows() map[string]string {
	return tukdbint.GetWorkflowDefinitionNames()
}
func GetActiveWorkflowNames() map[string]string {
	return tukdbint.GetActiveWorkflowNames()
}
func GetTaskNotes(pwy string, nhsid string, taskid int, ver int) string {
	return tukdbint.GetTaskNotes(pwy, nhsid, taskid, ver)
}
func (i *Transaction) IsTaskOverdue() bool {
	log.Printf("Checking if Workflow %s Task %v is overdue", i.Pathway, i.Task_ID)
	completionDate := i.GetTaskCompleteByDate()
	log.Printf("Task complete by time %s", completionDate)
	if time.Now().Local().Before(completionDate) {
		log.Printf("Time Now is before Task Complete by date. Task %v is NOT overdue", i.Task_ID)
		return false
	}
	if i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.TaskDetails.Status == tukcnst.COMPLETE {
		log.Printf("Task %v is Complete. Checking latest task event time", i.Task_ID)
		lasteventime := tukutil.GetTimeFromString(i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.TaskDetails.LastModifiedTime)
		if lasteventime.Before(completionDate) {
			log.Printf("Task %v was NOT overdue", i.Task_ID)
			return false
		}
	}
	log.Printf("Task %v IS overdue", i.Task_ID)
	return true
}
func (i *Transaction) GetTaskCompleteByDate() time.Time {
	if i.WorkflowDefinition.Tasks[i.Task_ID-1].CompleteByTime == "" {
		return i.GetWorkflowCompleteByDate()
	}
	return tukutil.OHT_FutureDate(tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value), i.WorkflowDefinition.Tasks[i.Task_ID-1].CompleteByTime)
}
func (i *Transaction) GetWorkflowDuration() string {
	ws := tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value)
	log.Printf("Workflow Started %s Status %s", ws.String(), i.WorkflowDocument.WorkflowStatus)
	loc, _ := time.LoadLocation("Europe/London")
	we := time.Now().In(loc)
	log.Printf("Time Now %s", we.String())
	if i.WorkflowDocument.WorkflowStatus == tukcnst.CLOSED {
		we = i.GetLatestWorkflowEventTime()
		log.Printf("Workflow is Complete. Latest Event Time was %s", we.String())
	}
	duration := tukutil.GetDuration(we.Sub(ws))
	log.Println("Duration - " + duration)
	return duration
}
func (i *Transaction) GetLatestWorkflowEventTime() time.Time {
	var we = tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value)

	for _, docevent := range i.WorkflowDocument.WorkflowStatusHistory.DocumentEvent {
		if docevent.EventTime != "" {
			doceventtime := tukutil.GetTimeFromString(docevent.EventTime)
			if doceventtime.After(we) {
				we = doceventtime
			}
		}
	}
	return we
}
func (i *Transaction) newEventID(eventType string, vers int) int64 {
	ev := tukdbint.Event{
		EventType:          eventType,
		DocName:            i.AttachmentInfo,
		ClassCode:          i.XDSDocumentMeta.Classcode,
		ConfCode:           i.XDSDocumentMeta.Confcode,
		FormatCode:         i.XDSDocumentMeta.Formatcode,
		FacilityCode:       i.XDSDocumentMeta.Facilitycode,
		PracticeCode:       i.XDSDocumentMeta.Practicesettingcode,
		Expression:         i.Expression,
		Authors:            i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Prefix + " " + i.WorkflowDocument.Author.AssignedAuthor.AssignedPerson.Name.Family,
		XdsPid:             i.XDS_ID,
		XdsDocEntryUid:     i.WorkflowDocument.ID.Extension,
		RepositoryUniqueId: i.XDSDocumentMeta.Repositoryuniqueid,
		NhsId:              i.NHS_ID,
		User:               i.User,
		Org:                i.Org,
		Role:               i.Role,
		Topic:              tukcnst.DSUB_TOPIC_TYPE_CODE,
		Pathway:            i.Pathway,
		Comments:           string(i.Request),
		Version:            vers,
		TaskId:             i.Task_ID,
	}
	evs := tukdbint.Events{Action: tukcnst.INSERT}
	evs.Events = append(evs.Events, ev)
	if err := tukdbint.NewDBEvent(&evs); err != nil {
		log.Println(err.Error())
		return 0
	}
	log.Printf("Created Event ID :  = %v", evs.LastInsertId)
	return evs.LastInsertId
}
func (i *Transaction) IsWorkflowOverdue() bool {
	if i.WorkflowDefinition.CompleteByTime != "" {
		completebyDate := i.GetWorkflowCompleteByDate()
		log.Printf("Workflow Complete By Date %s", completebyDate.String())
		if time.Now().After(completebyDate) {
			log.Printf("Time Now is after Workflow Complete By Date %s", completebyDate.String())
			if i.WorkflowDocument.WorkflowStatus == tukcnst.CLOSED {
				log.Println("Workflow is Complete, Obtaining latest workflow event time")
				levent := i.GetLatestWorkflowEventTime()
				log.Printf("Workflow Latest Event Time %s. Workflow Target Met = %v", levent.String(), levent.Before(completebyDate))
				return levent.After(completebyDate)
			} else {
				log.Printf("Workflow is not Complete. Complete By Date is %s Workflow Target not met", completebyDate.String())
				return true
			}
		} else {
			log.Printf("Time Now is before Workflow Complete By Date %s. Workflow is not overdue", completebyDate.String())
			return false
		}
	}
	log.Println("Workflow definition does not specify a Complete By Time. Workflow is not overdue")
	return false
}
func (i *Transaction) GetWorkflowCompleteByDate() time.Time {
	return tukutil.OHT_FutureDate(tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value), i.WorkflowDefinition.CompleteByTime)
}
func (i *Transaction) IsWorkflowCompleteBehaviorMet() (bool, string) {
	var conditions []string
	var completedConditions = 0
	var completionTask string
	for _, cc := range i.WorkflowDefinition.CompletionBehavior {
		if cc.Completion.Condition != "" {
			log.Println("Parsing Workflow Completion Condition " + cc.Completion.Condition)
			if strings.Contains(cc.Completion.Condition, " and ") {
				conditions = strings.Split(cc.Completion.Condition, " and ")
			} else {
				conditions = append(conditions, cc.Completion.Condition)
			}
			for _, condition := range conditions {
				endMethodInd := strings.Index(condition, "(")
				if endMethodInd > 0 {
					method := cc.Completion.Condition[0:endMethodInd]
					if method != "task" {
						log.Println(method + " is an Invalid Workflow Completion Behaviour Condition method. Ignoring Condition")
						continue
					}
					endParamInd := strings.Index(cc.Completion.Condition, ")")
					param := cc.Completion.Condition[endMethodInd+1 : endParamInd]
					for _, task := range i.WorkflowDocument.TaskList.XDWTask {
						if task.TaskData.TaskDetails.ID == param {
							if task.TaskData.TaskDetails.Status == tukcnst.COMPLETE {
								completedConditions = completedConditions + 1
								completionTask = param
							}
						}
					}
				}
			}
		}
	}
	if len(conditions) == completedConditions {
		log.Printf("%s Workflow for NHS ID %s is complete", i.Pathway, i.NHS_ID)
		return true, completionTask
	}
	log.Printf("%s Workflow for NHS ID %s is not complete", i.Pathway, i.NHS_ID)
	return false, ""
}
func IsLatestTaskEvent(i WorkflowDocument, task int, taskEventName string) bool {
	var lastoutputtime = tukutil.GetTimeFromString(i.EffectiveTime.Value)
	var lastinputtime = lastoutputtime
	var latestInputTaskEvent string
	var latestOutputTaskEvent string
	for _, in := range i.TaskList.XDWTask[task].TaskData.Input {
		if in.Part.AttachmentInfo.AttachedTime != "" {
			inputtime := tukutil.GetTimeFromString(in.Part.AttachmentInfo.AttachedTime)
			if inputtime.After(lastinputtime) {
				lastinputtime = inputtime
				latestInputTaskEvent = in.Part.AttachmentInfo.Name
			}
		}
	}
	for _, op := range i.TaskList.XDWTask[task].TaskData.Output {
		if op.Part.AttachmentInfo.AttachedTime != "" {
			outputtime := tukutil.GetTimeFromString(op.Part.AttachmentInfo.AttachedTime)
			if outputtime.After(lastoutputtime) {
				lastoutputtime = outputtime
				latestOutputTaskEvent = op.Part.AttachmentInfo.Name
			}
		}
	}
	if lastoutputtime.After(lastinputtime) {
		if taskEventName == latestOutputTaskEvent {
			return true
		}
	} else {
		if taskEventName == latestInputTaskEvent {
			return true
		}
	}
	return false
}
func (i *Transaction) IsTaskCompleteBehaviorMet() bool {
	log.Printf("Checking if Task %v is complete", i.Task_ID)
	var conditions []string
	var completedConditions = 0
	for _, cond := range i.WorkflowDefinition.Tasks[i.Task_ID].CompletionBehavior {
		log.Printf("Task %v Completion Condition is %s", i.Task_ID, cond)
		if cond.Completion.Condition != "" {
			if strings.Contains(cond.Completion.Condition, " and ") {
				conditions = strings.Split(cond.Completion.Condition, " and ")
			} else {
				conditions = append(conditions, cond.Completion.Condition)
			}
			log.Printf("Checkiing Task %v %v completion conditions", i.Task_ID, len(conditions))
			for _, condition := range conditions {
				endMethodInd := strings.Index(condition, "(")
				if endMethodInd > 0 {
					method := condition[0:endMethodInd]
					endParamInd := strings.Index(condition, ")")
					if endParamInd < endMethodInd+2 {
						log.Println("Invalid Condition. End bracket index invalid")
						continue
					}
					param := condition[endMethodInd+1 : endParamInd]
					log.Printf("Completion condition is %s", method)
					switch method {
					case "output":
						for _, op := range i.WorkflowDocument.TaskList.XDWTask[i.Task_ID].TaskData.Output {
							if op.Part.AttachmentInfo.AttachedTime != "" && op.Part.AttachmentInfo.Name == param {
								completedConditions = completedConditions + 1
							}
						}
					case "input":
						for _, in := range i.WorkflowDocument.TaskList.XDWTask[i.Task_ID].TaskData.Input {
							if in.Part.AttachmentInfo.AttachedTime != "" && in.Part.AttachmentInfo.Name == param {
								completedConditions = completedConditions + 1
							}
						}
					case "task":
						for _, task := range i.WorkflowDocument.TaskList.XDWTask {
							if task.TaskData.TaskDetails.ID == param {
								if task.TaskData.TaskDetails.Status == tukcnst.COMPLETE {
									completedConditions = completedConditions + 1
								}
							}
						}
					case "latest":
						if i.getLatestTaskEvent() == param {
							completedConditions = completedConditions + 1
						}
					}
				}
			}
		}
	}
	if len(conditions) == completedConditions {
		log.Printf("Task %v is complete", i.Task_ID)
		return true
	}
	log.Printf("Task %v is not complete", i.Task_ID)
	return false
}
func (i *Transaction) getLatestTaskEvent() string {
	var lasteventtime = tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value)
	var lastevent = ""
	for _, v := range i.WorkflowDocument.TaskList.XDWTask[i.Task_ID].TaskData.Input {
		if v.Part.AttachmentInfo.AttachedTime != "" {
			et := tukutil.GetTimeFromString(v.Part.AttachmentInfo.AttachedTime)
			if et.After(lasteventtime) {
				lasteventtime = et
				lastevent = v.Part.AttachmentInfo.Name
			}
		}
	}
	for _, v := range i.WorkflowDocument.TaskList.XDWTask[i.Task_ID].TaskData.Output {
		if v.Part.AttachmentInfo.AttachedTime != "" {
			et := tukutil.GetTimeFromString(v.Part.AttachmentInfo.AttachedTime)
			if et.After(lasteventtime) {
				lasteventtime = et
				lastevent = v.Part.AttachmentInfo.Name
			}
		}
	}
	return lastevent
}
func (i *Transaction) GetTaskDuration() string {
	taskCreationTime := tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value)
	log.Printf("Task %v Creation Time %s", i.Task_ID, taskCreationTime.String())
	if i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.TaskDetails.Status == tukcnst.COMPLETE {
		log.Printf("Workflow Task %s is complete", i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.TaskDetails.Name)
		lastEvent := tukutil.GetTimeFromString(i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.TaskDetails.LastModifiedTime)
		log.Printf("Lastest Task Event %s", lastEvent.String())
		duration := lastEvent.Sub(taskCreationTime)
		log.Printf("Task %v %s Created %s Status is COMPLETE Duration - %s", i.Task_ID, i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.Description, taskCreationTime.String(), duration.String())
		return tukutil.PrettyPrintDuration(duration)
	} else {
		duration := time.Since(taskCreationTime)
		log.Printf("Task %v %s Created %s Status is %s Duration - %s", i.Task_ID, i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.Description, taskCreationTime.String(), i.WorkflowDocument.TaskList.XDWTask[i.Task_ID-1].TaskData.TaskDetails.Status, duration.String())
		return tukutil.PrettyPrintDuration(duration)
	}
}
func (i *Transaction) GetTaskTimeRemaining() string {
	taskCreateTime := tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value)
	taskCompleteby := tukutil.OHT_FutureDate(taskCreateTime, i.WorkflowDefinition.Tasks[i.Task_ID-1].CompleteByTime)
	log.Printf("Completion time %s", taskCompleteby.String())
	if time.Now().After(taskCompleteby) {
		return "0"
	}
	timeRemaining := taskCompleteby.Sub(taskCreateTime)
	log.Println("Task Time Remaining : " + timeRemaining.String())
	return tukutil.PrettyPrintDuration(timeRemaining)
}
func (i *Transaction) GetWorkflowTimeRemaining() string {
	createTime := tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value)
	completeby := tukutil.OHT_FutureDate(createTime, i.WorkflowDefinition.CompleteByTime)
	log.Printf("Completion time %s", completeby.String())
	if time.Now().After(completeby) {
		return "0"
	}
	timeRemaining := time.Until(completeby)
	log.Println("Workflow Time Remaining : " + timeRemaining.String())
	return tukutil.PrettyPrintDuration(timeRemaining)
}
func getLocalId(mid string) string {
	return tukdbint.GetIDMapsLocalId(mid)
}
func (i *Transaction) IsWorkflowEscalated() bool {
	if i.WorkflowDefinition.ExpirationTime != "" {
		escalatedate := tukutil.OHT_FutureDate(tukutil.GetTimeFromString(i.WorkflowDocument.EffectiveTime.Value), i.WorkflowDefinition.ExpirationTime)
		log.Printf("Workflow Start Time %s Worklow Escalate Time %s Workflow Escaleted = %v", i.WorkflowDocument.EffectiveTime.Value, escalatedate.String(), time.Now().After(escalatedate))
		return time.Now().After(escalatedate)
	}
	log.Println("No Escalate time defined for Workflow")
	return false
}

// XDW Admin functions

func (i *Transaction) RegisterWorkflowDefinition(ismeta bool) error {
	var err error
	if i.Pathway == "" {
		return errors.New("pathway is not set")
	}
	if i.Request == nil || string(i.Request) == "" {
		return errors.New("request bytes is not set")
	}
	if ismeta {
		log.Println("Persisting Workflow XDS Meta")
		err = i.registerWorkflowXDSMeta()
	} else {
		log.Println("Registering Workflow Definition")
		err = i.registerWorkflowDef()
	}
	return err
}
func (i *Transaction) registerWorkflowDef() error {
	var err error
	err = json.Unmarshal(i.Request, &i.WorkflowDefinition)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	event := tukdsub.DSUBEvent{Action: tukcnst.CANCEL, Pathway: i.WorkflowDefinition.Ref}
	tukdsub.New_Transaction(&event)
	log.Printf("Cancelled any existing DSUB Broker and Event Service Subscriptions for Pathway %s", i.WorkflowDefinition.Ref)
	dsubSubs := make(map[string]string)
	if err = i.PersistXDWDefinition(); err == nil {
		log.Println("Parsing XDW Tasks for potential DSUB Broker Subscriptions")
		for _, task := range i.WorkflowDefinition.Tasks {
			for _, inp := range task.Input {
				log.Printf("Checking Task %v %s input %s", task.ID, task.Name, inp.Name)
				switch inp.AccessType {
				case tukcnst.XDS_REGISTERED:
					dsubSubs[inp.Name] = i.WorkflowDefinition.Ref
					log.Printf("Task %v %s input %s included in potential DSUB Broker subscriptions", task.ID, task.Name, inp.Name)
				}
			}
			for _, out := range task.Output {
				log.Printf("Checking Task %v %s output %s", task.ID, task.Name, out.Name)
				switch out.AccessType {
				case tukcnst.XDS_REGISTERED:
					dsubSubs[out.Name] = i.WorkflowDefinition.Ref
					log.Printf("Task %v %s input %s included in potential DSUB Broker subscriptions", task.ID, task.Name, out.Name)
				}
			}
		}
	}
	log.Printf("Found %v potential DSUB Broker Subscriptions - %s", len(dsubSubs), dsubSubs)
	if len(dsubSubs) > 0 {
		event.Action = tukcnst.CREATE
		event.BrokerURL = i.ServiceURL.DSUB_Broker
		event.ConsumerURL = i.ServiceURL.DSUB_Consumer
		for expression := range dsubSubs {
			event.Expressions = append(event.Expressions, expression)
		}
		err = tukdsub.New_Transaction(&event)
	}

	return err
}
func (i *Transaction) registerWorkflowXDSMeta() error {
	var err error
	xdw := tukdbint.XDW{Name: i.Pathway + "_meta", IsXDSMeta: true}
	xdws := tukdbint.XDWS{Action: tukcnst.DELETE}
	xdws.XDW = append(xdws.XDW, xdw)
	if err = tukdbint.NewDBEvent(&xdws); err == nil {
		log.Printf("Deleted Existing XDS Meta for Pathway %s", i.Pathway)
		xdw = tukdbint.XDW{Name: i.Pathway + "_meta", IsXDSMeta: true, XDW: string(i.Request)}
		xdws = tukdbint.XDWS{Action: tukcnst.INSERT}
		xdws.XDW = append(xdws.XDW, xdw)
		if err = tukdbint.NewDBEvent(&xdws); err == nil {
			log.Printf("Persisted Workflow XDS Meta for Pathway %s", i.Pathway)
		}
	}
	return err
}
func (i *Transaction) PersistXDWDefinition() error {
	var err error
	xdw := tukdbint.XDW{Name: i.Pathway, IsXDSMeta: false}
	xdws := tukdbint.XDWS{Action: tukcnst.DELETE}
	xdws.XDW = append(xdws.XDW, xdw)
	if err = tukdbint.NewDBEvent(&xdws); err == nil {
		xdw = tukdbint.XDW{Name: i.Pathway, IsXDSMeta: false, XDW: string(i.Request)}
		xdws = tukdbint.XDWS{Action: tukcnst.INSERT}
		xdws.XDW = append(xdws.XDW, xdw)
		if err = tukdbint.NewDBEvent(&xdws); err == nil {
			log.Printf("Persisted New XDW Definition for Pathway %s", i.Pathway)
		}
	}
	return err
}
func GetWorkflowDefinitionNames() []string {
	var wfnames []string
	names := tukdbint.GetWorkflowDefinitionNames()
	for name := range names {
		wfnames = append(wfnames, name)
	}
	return wfnames
}
func GetWorkflowXDSMetaNames() []string {
	return tukdbint.GetWorkflowXDSMetaNames()
}

// sort interface
type DocumentEvents []DocumentEvent

func (e DocumentEvents) Len() int {
	return len(e)
}
func (e DocumentEvents) Less(i, j int) bool {
	return e[i].EventTime > e[j].EventTime
}
func (e DocumentEvents) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type eventsList []tukdbint.Event

func (e eventsList) Len() int {
	return len(e)
}
func (e eventsList) Less(i, j int) bool {
	return e[i].Id > e[j].Id
}
func (e eventsList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
