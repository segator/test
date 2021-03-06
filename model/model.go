package model

import (
	"github.com/google/uuid"
	"time"
	"transcoder/helper/max"
)
type EventType string
type NotificationType string
type NotificationStatus string
type JobAction string
type TaskEvents []*TaskEvent
const(
	PingEvent         EventType = "Ping"
	NotificationEvent EventType = "Notification"

	JobNotification        NotificationType = "Job"
	DownloadNotification   NotificationType = "Download"
	UploadNotification   NotificationType = "Upload"
	MKVExtractNotification NotificationType = "MKVExtract"
	FFProbeNotification    NotificationType = "FFProbe"
	PGSNotification        NotificationType = "PGS"
	FFMPEGSNotification    NotificationType = "FFMPEG"

	AddedNotificationStatus     NotificationStatus = "added"
	ReAddedNotificationStatus   NotificationStatus = "readded"
	StartedNotificationStatus   NotificationStatus = "started"
	CompletedNotificationStatus NotificationStatus = "completed"
	CanceledNotificationStatus  NotificationStatus = "canceled"
	FailedNotificationStatus    NotificationStatus = "failed"

	CancelJob JobAction = "cancel"
	EncodeJobType JobType = "encode"
	PGSToSrtJobType JobType = "pgstosrt"
)
type Identity interface{
	getUUID() uuid.UUID
}
type Video struct {
	SourcePath      string     `json:"sourcePath"`
	DestinationPath string     `json:"destinationPath"`
	Id   uuid.UUID `json:"id"`
	Events          TaskEvents `json:"events"`
}

type JobEventQueue struct {
	Queue string
	JobEvent *JobEvent
}
type Worker struct {
	Name string
	Ip string
	QueueName string
	LastSeen time.Time
}

type ControlEvent struct {
	Event *TaskEncode
	ControlChan chan interface{}
}

type JobEvent struct {
	Id   uuid.UUID `json:"id"`
	Action JobAction `json:"action"`
}

type JobType string

type TaskEncode struct {
	Id   uuid.UUID `json:"id"`
	DownloadURL string    `json:"downloadURL"`
	UploadURL   string    `json:"uploadURL"`
	ChecksumURL string    `json:"checksumURL"`
	EventID     int `json:"eventID"`
	Priority    int `json:"priority"`
}

type TaskPGS struct {
	Id   uuid.UUID `json:"id"`
	PGSID  int   `json:"pgsid"`
	PGSdata  []byte    `json:"pgsdata"`
	ReplyTo string `json:"replyto"`
}

type TaskPGSResponse struct {
	Id   uuid.UUID `json:"id"`
	PGSID  int   `json:"pgsid"`
	Srt []byte `json:"srt"`
	Err string `json:"error"`
}


func (V TaskEncode) getUUID() uuid.UUID{
	return V.Id
}
func (V TaskPGS) getUUID() uuid.UUID{
	return V.Id
}


type TaskEvent struct {
	Id   uuid.UUID `json:"id"`
	EventID          int                `json:"eventID"`
    EventType        EventType          `json:"eventType"`
	WorkerName       string             `json:"workerName"`
	WorkerQueue      string             `json:"workerQueue"`
	EventTime        time.Time          `json:"eventTime"`
	IP               string             `json:"ip"`
	NotificationType NotificationType   `json:"notificationType"`
	Status           NotificationStatus `json:"status"`
	Message          string             `json:"message"`
}

func (t *TaskEvents) GetLatest() *TaskEvent {
	if len(*t)==0{
		return nil
	}
	return max.Max(t).(*TaskEvent)
}
func (t *TaskEvents) GetLatestPerNotificationType(notificationType NotificationType) (returnEvent *TaskEvent){
	eventID:=-1
	for _,event :=range *t {
		if event.NotificationType == notificationType && event.EventID>eventID{
			eventID=event.EventID
			returnEvent=event
		}
	}
	return returnEvent
}
func (t *TaskEvents) GetStatus() NotificationStatus {
	return t.GetLatestPerNotificationType(JobNotification).Status
}

type JobRequestError struct {
	JobRequest
	Error string `json:"error"`
}
type JobRequest struct {
	SourcePath string `json:"sourcePath"`
	DestinationPath string `json:"destinationPath"`
	ForceCompleted  bool `json:"forceCompleted"`
	ForceFailed  bool `json:"forceFailed"`
	ForceExecuting  bool `json:"forceExecuting"`
	ForceAdded bool `json:"forceAdded"`
	Priority int `json:"priority"`
}



func (a TaskEvents) Len() int{
	return len(a)
}
func (a TaskEvents) Less(i, j int) bool {
	return a[i].EventID <a[j].EventID
}
func (a TaskEvents) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a TaskEvents) GetLastElement(i int) interface{}{
	return a[i]
}

func (v *Video) AddEvent(eventType EventType,notificationType NotificationType, notificationStatus NotificationStatus) (newEvent *TaskEvent) {
	latestEvent := v.Events.GetLatest()
	newEventID:=0
	if latestEvent!=nil{
		newEventID=latestEvent.EventID+1
	}

	newEvent = &TaskEvent{
		Id:         v.Id,
		EventID:          newEventID,
		EventType:        eventType,
		EventTime:        time.Now(),
		NotificationType: notificationType,
		Status:           notificationStatus,
	}
	v.Events=append(v.Events,newEvent)
	return newEvent
}

type WorkerQueue interface{
	RegisterWorker(worker QueueWorker)
}
type QueueWorker interface{
	IsTypeAccepted(jobType string) bool
	Prepare(workData[]byte,queueManager Manager) error
	Execute() error
	Clean() error
	Cancel()
	GetID() string
	GetTaskID() uuid.UUID
}
type Manager interface{
	EventNotification(event TaskEvent)
	ResponsePGSJob(response TaskPGSResponse) error
	RequestPGSJob(pgsJob TaskPGS) <-chan TaskPGSResponse
}
