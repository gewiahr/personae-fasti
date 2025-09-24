package reqData

type RecordInsert struct {
	Text     string `json:"text"`
	Hidden   bool   `json:"hidden"`
	PlayerID int    `json:"-"`
	GameID   int    `json:"-"`
	QuestID  int    `json:"questID"`
}

type RecordUpdate struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Hidden  bool   `json:"hidden"`
	QuestID int    `json:"questID"`
}

type CharCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

type CharUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

type NPCCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

type NPCUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

type LocationCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ParentID    int    `json:"pid"`
	Hidden      bool   `json:"hidden"`
}

type LocationUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ParentID    int    `json:"pid"`
	Hidden      bool   `json:"hidden"`
}

type QuestCreateData struct {
	Quest QuestCreate  `json:"quest"`
	Tasks []TaskCreate `json:"tasks"`
}

type QuestUpdateData struct {
	Quest QuestUpdate  `json:"quest"`
	Tasks []TaskUpdate `json:"tasks"`
}

type QuestTasksPatch struct {
	QuestID int         `json:"questID"`
	Tasks   []TaskPatch `json:"tasks"`
}

type QuestCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	ParentID int `json:"parentID"`
	ChildID  int `json:"childID"`
	HeadID   int `json:"headID"`

	Successful bool `json:"successful"`

	Hidden bool `json:"hidden"`
}

type QuestUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	ParentID int `json:"parentID"`
	ChildID  int `json:"childID"`
	HeadID   int `json:"headID"`

	Successful bool `json:"successful"`

	Hidden bool `json:"hidden"`

	Finished bool `json:"finished"`
}

type TaskCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Capacity    int    `json:"capacity"`

	Hidden bool `json:"hidden"`
}

type TaskUpdate struct {
	ID int `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Capacity    int    `json:"capacity"`

	Hidden bool `json:"hidden"`
}

type TaskPatch struct {
	ID      int `json:"id"`
	Current int `json:"current"`
}

type GameChange struct {
	GameID int `json:"gameID"`
}

type GameSettingsUpdate struct {
	GameID              int  `json:"gameID"`
	AllowAllEditRecords bool `json:"allowAllEditRecords"`
}
