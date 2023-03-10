package bot

import (
	"fmt"
	"github.com/artchitector/artchitect/model"
)

func getTextWithoutCaption(card model.Card) string {
	return fmt.Sprintf(
		"Card #%d. (https://artchitect.space/card/%d)\n\n"+
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

func getTextWithCaption(card model.Card, caption string) string {
	return fmt.Sprintf(
		"\"%s\"\n"+
			"Card #%d. (https://artchitect.space/card/%d)\n\n"+
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

func getUnityText(unity model.Unity) string {
	return fmt.Sprintf(
		"Unity %s unified. Version %d\n https://artchitect.space/unity/%s",
		unity.Mask,
		unity.Version,
		unity.Mask,
	)
}
