#!/bin/bash
# verify_ko.sh — verifies that the Korean (ko) localization loads correctly
# and returns Korean strings for all keys that have Korean translations.
#
# Usage: ./verify_ko.sh [path/to/end_of_eden]
#
# Requires: go (any 1.21+ version)
# Run from the repo root (or pass the path as the first argument).

set -e

REPO_ROOT="${1:-$(pwd)}"
if [ ! -d "$REPO_ROOT/assets/locals/ko" ]; then
    echo "ERROR: $REPO_ROOT/assets/locals/ko not found."
    echo "Run this script from the end_of_eden repo root or pass the path."
    exit 1
fi

WORK=$(mktemp -d)
trap "rm -rf $WORK" EXIT

cat > "$WORK/main.go" <<EOF
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BigJk/end_of_eden/system/localization"
)

func main() {
	if err := localization.Global.AddFolder("$REPO_ROOT/assets/locals"); err != nil {
		fmt.Println("AddFolder error:", err)
		os.Exit(1)
	}

	koKeys := []string{
		"intents.deal_damage","intents.block","intents.heal_ally",
		"intents.charge_laser","intents.charge_plasma","intents.wait",
		"intents.load_battery","intents.standby","intents.deal_heal",
		"events.START.name","events.START.description",
		"events.RUST_MITE.name","events.RUST_MITE.description","events.RUST_MITE.choices.0.description",
		"events.MERCHANT.name","events.MERCHANT.description","events.MERCHANT.choices.0.description","events.MERCHANT.choices.1.description",
		"events.ARTIFACT_CHOICE.name","events.ARTIFACT_CHOICE.description","events.ARTIFACT_CHOICE.take_template","events.ARTIFACT_CHOICE.skip",
		"events.CARD_CHOICE.name","events.CARD_CHOICE.description","events.CARD_CHOICE.take_template","events.CARD_CHOICE.skip",
		"events.RANDOM_ARTIFACT_ACT_0.name","events.RANDOM_ARTIFACT_ACT_0.description","events.RANDOM_ARTIFACT_ACT_0.choices.1.description",
		"events.RANDOM_CONSUMEABLE_ACT_0.name","events.RANDOM_CONSUMEABLE_ACT_0.choices.1.description",
		"events.GAIN_GOLD_ACT_0.name","events.GAIN_GOLD_ACT_0.description","events.GAIN_GOLD_ACT_0.choices.1.description",
		"events.GOLD_TO_HP_ACT_0.name","events.GOLD_TO_HP_ACT_0.choices.1.description",
		"events.MAX_LIFE_ACT_0.name","events.MAX_LIFE_ACT_0.choices.1.description",
		"events.GAMBLE_1_ACT_0.name","events.GAMBLE_1_ACT_0.choices.1.description",
		"events.UPRAGDE_CARD_ACT_0.name","events.UPRAGDE_CARD_ACT_0.choices.1.description",
		"ui.menu.title","ui.menu.continue","ui.menu.continue_desc",
		"ui.menu.tutorial","ui.menu.tutorial_desc","ui.menu.new_game","ui.menu.new_game_desc",
		"ui.menu.new_game_sod","ui.menu.new_game_sod_desc","ui.menu.about","ui.menu.about_desc",
		"ui.menu.settings","ui.menu.settings_desc","ui.menu.mods","ui.menu.mods_desc",
		"ui.menu.exit","ui.menu.exit_desc",
		"ui.overview.character","ui.overview.character_desc","ui.overview.logs","ui.overview.logs_desc",
		"ui.overview.artifacts","ui.overview.artifacts_desc","ui.overview.cards","ui.overview.cards_desc",
		"ui.overview.quit","ui.overview.quit_desc",
		"ui.merchant.wares","ui.merchant.leave","ui.merchant.buy_item","ui.merchant.upgrade_card_fmt",
		"ui.gameview.end_turn","ui.gameview.intend","ui.gameview.status_effects","ui.gameview.player_status","ui.gameview.press_esc",
		"ui.gameover.run_statistic","ui.gameover.stages","ui.gameover.damage_done","ui.gameover.damage_received","ui.gameover.gold_collected","ui.gameover.accept_fate",
		"ui.carousel.continue",
		"ui.lua_error.back","ui.lua_error.copied","ui.lua_error.copy_clipboard","ui.lua_error.title","ui.lua_error.error",
		"ui.settings.title",
		"ui.about_text","ui.about.version_fmt","ui.about.back",
		"log.type.info","log.type.warning","log.type.danger","log.type.success",
		"log.session.saving","log.player.hit_enemy","log.player.took_damage",
		"log.enemy.took_damage_fmt","log.enemy.died_dropped_gold",
		"log.player.healed","log.enemy.healed",
		"enemies.TUTORIAL_DUMMY_1.name","enemies.TUTORIAL_DUMMY_1.description",
		"enemies.TUTORIAL_DUMMY_2.name","enemies.TUTORIAL_DUMMY_2.description",
		"status_effects.WEAKNESS.name","status_effects.WEAKNESS.description","status_effects.WEAKNESS.state",
		"events.TUTORIAL_START.name","events.TUTORIAL_START.description","events.TUTORIAL_START.choices.0.description",
		"events.TUTORIAL_1.name","events.TUTORIAL_1.description","events.TUTORIAL_1.choices.0.description",
		"events.TUTORIAL_2.name","events.TUTORIAL_2.description","events.TUTORIAL_2.choices.0.description","events.TUTORIAL_2.intent",
		"events.TUTORIAL_3.name","events.TUTORIAL_3.description","events.TUTORIAL_3.choices.0.description",
		"events.TUTORIAL_4.name","events.TUTORIAL_4.description","events.TUTORIAL_4.choices.0.description",
		"enemies.CLEAN_BOT.name","enemies.LASER_DRONE.name","enemies.PLASMA_GOLEM.name","enemies.CYBER_SPIDER.name",
		"enemies.CYBER_SLIME_MINION.name","enemies.CYBER_SLIME.name","enemies.REPAIR_DRONE.name",
		"status_effects.CHARGING.description","status_effects.CHARGING.state",
		"artifacts.HEADBUT_HELMET.name","artifacts.HEADBUT_HELMET.description",
		"ui.one_time","ui.heal_2",
		"ui.found_template","ui.found_event_template","ui.take_template","ui.leave",
		"artifacts.INTERVA_JUICER.description_suffix",
		"weapons.crowbar.name","weapons.crowbar.description",
		"weapons.vibro_knife.name","weapons.vibro_knife.description",
		"weapons.lzr_pistol.name","weapons.lzr_pistol.description",
		"weapons.har_ii.name","weapons.har_ii.description",
		"cards.SMOKE_BOMB.name","cards.SMOKE_BOMB.description",
		"status_effects.SMOKE_BOMB.name","status_effects.SMOKE_BOMB.description",

		"basics.on","basics.off",
		"settings.audio.title","settings.volume.description",
		"settings.language.title","settings.crt.title","settings.fps.title",
		"cards.MELEE_HIT.name","cards.MELEE_HIT.description","cards.MELEE_HIT.state",
		"cards.ARM_MOUNTED_GUN.name","cards.ARM_MOUNTED_GUN.description","cards.ARM_MOUNTED_GUN.state",
		"cards.CROWBAR.name","cards.VIBRO_KNIFE.name","cards.LZR_PISTOL.name","cards.HAR-II.name",
		"cards.ADRENALINE_SHOT.name","cards.ADRENALINE_SHOT.description",
		"cards.ENRGY_DRINK_X91.name","cards.ENRGY_DRINK_X92.name","cards.ENRGY_DRINK_X93.name",
		"cards.BLOCK.name","cards.BLOCK.description","cards.BLOCK.state",
		"cards.BOUNCE_SHIELD.name","cards.BOUNCE_SHIELD.description","cards.BOUNCE_SHIELD.state",
		"cards.FLASH_BANG.name","cards.FLASH_BANG.description",
		"cards.FLASH_SHIELD.name","cards.FLASH_SHIELD.description","cards.FLASH_SHIELD.state",
		"cards.KNOCK_OUT.name","cards.KNOCK_OUT.description",
		"cards.LONG_REST.name","cards.LONG_REST.description","cards.LONG_REST.state",
		"cards.NANO_CHARGER.name","cards.NANO_CHARGER.description","cards.NANO_CHARGER.state",
		"cards.NULLIFY.name","cards.NULLIFY.description","cards.NULLIFY.state",
		"cards.REST.name","cards.REST.description","cards.REST.state",
		"cards.SMOKE_BOMB.name","cards.SMOKE_BOMB.description","cards.SMOKE_BOMB.state",
		"cards.STIM_PACK.name","cards.STIM_PACK.description","cards.STIM_PACK.state",
		"cards.ULTRA_FLASH_SHIELD.name","cards.ULTRA_FLASH_SHIELD.description","cards.ULTRA_FLASH_SHIELD.state",
		"status_effects.ADRENALINE_SHOT.name","status_effects.ADRENALINE_SHOT.description",
		"status_effects.BLOCK.name","status_effects.BLOCK.description",
		"status_effects.BOUNCE_SHIELD.name","status_effects.BOUNCE_SHIELD.description",
		"status_effects.FLASH_BANG.name","status_effects.FLASH_BANG.description",
		"status_effects.FLASH_SHIELD.name","status_effects.FLASH_SHIELD.description",
		"status_effects.KNOCK_OUT.name","status_effects.KNOCK_OUT.description",
		"status_effects.NANO_CHARGER.name","status_effects.NANO_CHARGER.description",
		"status_effects.NULLIFY.name","status_effects.NULLIFY.description",
		"status_effects.SMOKE_BOMB.name","status_effects.SMOKE_BOMB.description",
		"status_effects.ULTRA_FLASH_SHIELD.name","status_effects.ULTRA_FLASH_SHIELD.description",
		"status_effects.CHARGED.name","status_effects.CHARGED.description",
		"status_effects.CHARGING.name","status_effects.CHARGING.description",
		"artifacts.ARM_MOUNTED_GUN.name","artifacts.ARM_MOUNTED_GUN.description",
		"artifacts.BIO_RECYCLER.name","artifacts.BIO_RECYCLER.description",
		"artifacts.COMBAT_GLASSES.name","artifacts.COMBAT_GLASSES.description",
		"artifacts.COMBAT_GLOVES.name","artifacts.COMBAT_GLOVES.description",
		"artifacts.GOLD_SCRAPPER.name","artifacts.GOLD_SCRAPPER.description",
		"artifacts.HEADBUT_HELMET.name","artifacts.HEADBUT_HELMET.description",
		"artifacts.INTERVA_JUICER.name","artifacts.INTERVA_JUICER.description",
		"artifacts.PORTABLE_BUFFER.name","artifacts.PORTABLE_BUFFER.description",
		"artifacts.REFLECTIVE_ARMOR.name","artifacts.REFLECTIVE_ARMOR.description",
		"artifacts.SPEED_ENHANCER.name","artifacts.SPEED_ENHANCER.description",
		"enemies.RUST_MITE.name","enemies.RUST_MITE.description",
		"enemies.CLEAN_BOT.name","enemies.CLEAN_BOT.description",
		"enemies.LASER_DRONE.name","enemies.LASER_DRONE.description",
		"enemies.NANOBOT_SWARM.name","enemies.NANOBOT_SWARM.description",
		"enemies.CYBER_SLIME.name","enemies.CYBER_SLIME.description",
		"enemies.CYBER_SLIME_MINION.name","enemies.CYBER_SLIME_MINION.description",
		"enemies.CYBER_SPIDER.name","enemies.CYBER_SPIDER.description",
		"enemies.PLASMA_GOLEM.name","enemies.PLASMA_GOLEM.description",
		"enemies.REPAIR_DRONE.name","enemies.REPAIR_DRONE.description",
	}

	pass, fail := 0, 0
	var missing []string
	for _, k := range koKeys {
		v := localization.Global.Get("ko", k)
		if v == k || strings.HasPrefix(v, k) {
			fail++
			missing = append(missing, k)
		} else {
			pass++
		}
	}
	fmt.Printf("\n========== Korean (ko) localization check ==========\n")
	fmt.Printf("PASS: %d\nFAIL: %d\nTOTAL: %d\n", pass, fail, len(koKeys))
	if fail > 0 {
		fmt.Println("\nMissing keys (fell back to en or returned the key):")
		for _, k := range missing {
			fmt.Println("  -", k)
		}
		os.Exit(1)
	}
	fmt.Println("\nAll Korean keys load correctly.")
}
EOF

cd "$REPO_ROOT"
go run "$WORK/main.go"