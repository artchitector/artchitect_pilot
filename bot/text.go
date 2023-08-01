package bot

import (
	"encoding/json"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
)

func getTextWithoutCaption(card model.Art) string {
	return fmt.Sprintf(
		"Art #%d. (https://artchitect.space/card/%d)\n\n"+
			"Created: %s\n"+
			"Seed: %d\n"+
			"Tags: %s",
		card.ID,
		card.ID,
		card.CreatedAt.Format("2006 Jan 2 15:04"),
		card.Spell.Seed,
		card.Spell.Tags,
	)
}

func getTextWithCaption(card model.Art, caption string) string {
	return fmt.Sprintf(
		"\"%s\"\n"+
			"Art #%d. (https://artchitect.space/card/%d)\n\n"+
			"Created: %s\n"+
			"Seed: %d\n"+
			"Tags: %s",
		caption,
		card.ID,
		card.ID,
		card.CreatedAt.Format("2006 Jan 2 15:04"),
		card.Spell.Seed,
		card.Spell.Tags,
	)
}

func getUnityText(unity model.Unity) (string, error) {
	var leads []uint
	if err := json.Unmarshal([]byte(unity.Leads), &leads); err != nil {
		return "", errors.Wrapf(err, "[unifier] failed to unmarshal leads %s", unity.Mask)
	}
	totalLeads := len(leads)
	var filledLeads int
	for _, lead := range leads {
		if lead > 0 {
			filledLeads += 1

		}
	}
	isCompleted := totalLeads == filledLeads
	return fmt.Sprintf(
		"Unity %s unified. Version %d.\nLeads: %d from %d. Completed: %t.\n\nhttps://artchitect.space/unity/%s",
		unity.Mask,
		unity.Version,
		filledLeads,
		totalLeads,
		isCompleted,
		unity.Mask,
	), nil
}
