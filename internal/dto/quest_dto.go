package dto

/* Quest */

type QuestBrief struct {
	ExtID string `json:"ext"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameExtID  string `json:"gameExt"`
	Hidden     bool   `json:"hidden"`
	Successful bool   `json:"successful"`
	Finished   bool   `json:"finished"`
}

type QuestFull struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	GameExtID  string `json:"gameExt"`
	Hidden     bool   `json:"hidden"`
	Successful bool   `json:"successful"`
	Finished   bool   `json:"finished"`
}

type QuestPage struct {
	Quest   QuestFull       `json:"quest"`
	Tasks   []QuestTaskFull `json:"tasks"`
	Records []RecordFull    `json:"records"`
}

// type GameQuests struct {
// 	Quests []QuestInfo `json:"quests"`
// 	//CurrentGame GameInfo    `json:"currentGame"`
// }

/* Task */

type QuestTaskFull struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Type     int  `json:"type"`
	Capacity int  `json:"capacity"`
	Current  int  `json:"current"`
	Finished bool `json:"finished"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

type QuestTasksPatch struct {
	QuestExtID string      `json:"questExt"`
	Tasks      []TaskPatch `json:"tasks"`
}

type TaskPatch struct {
	ID      int `json:"id"`
	Current int `json:"current"`
}

/* Create Request */

type QuestCreateData struct {
	Quest QuestCreate  `json:"quest"`
	Tasks []TaskCreate `json:"tasks"`
}

type QuestCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Successful bool `json:"successful"`

	Hidden bool `json:"hidden"`
}

type TaskCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Capacity    int    `json:"capacity"`

	Hidden bool `json:"hidden"`
}

/* Update Request */

type QuestUpdateData struct {
	Quest QuestUpdate  `json:"quest"`
	Tasks []TaskUpdate `json:"tasks"`
}

type QuestUpdate struct {
	ExtID       string `json:"ext"`
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Successful bool `json:"successful"`

	Hidden bool `json:"hidden"`

	Finished bool `json:"finished"`
}

type TaskUpdate struct {
	ID int `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Capacity    int    `json:"capacity"`

	Hidden bool `json:"hidden"`
}
