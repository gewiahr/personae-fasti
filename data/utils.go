package data

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
)

func (s *Storage) InsertMentionsForRecord(record *Record) error {
	re, err := regexp.Compile(`@(?P<type>\w+):(?P<id>\d+)` + "`(?P<name>[^`]+)`")
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(record.Text, -1)
	for _, match := range matches {
		// Parse mention ID
		id, err := strconv.Atoi(match[2]) //ParseInt(match[2], 10, 64)
		if err != nil {
			continue
		}
		// Insert to a correct type
		switch match[1] {
		case "char":
			_, err = s.db.NewInsert().Model(&RecordChar{RecordID: record.ID, CharID: id}).Exec(context.Background())
			break
		case "npc":
			_, err = s.db.NewInsert().Model(&RecordNPC{RecordID: record.ID, NPCID: id}).Exec(context.Background())
			break
		case "location":
			_, err = s.db.NewInsert().Model(&RecordLocation{RecordID: record.ID, LocationID: id}).Exec(context.Background())
			break
		default:
			fmt.Printf("error during record mention extracting: mention %s is incorrect in record %d", match[0], record.ID)
			// add error logger
			break
		}
		// Return on Insert Error
		if err != nil {
			return err
		}
	}

	return err
}
