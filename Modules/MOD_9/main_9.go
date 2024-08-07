/*******************************************************************************
 * Copyright 2023-2024 Edw590
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package MOD_9

import (
	"RRComm/RRComm"
	"Registry/Registry"
	MOD_3 "Speech"
	"SpeechQueue/SpeechQueue"
	"ULComm/ULComm"
	MOD_12 "UserLocator"
	"Utils"
	"VISOR_Client/ClientRegKeys"
	"bytes"
	"github.com/apaxa-go/eval"
	"log"
	"strconv"
	"strings"
	"time"
)

// Reminders Reminder //

const TIME_SLEEP_S int = 1

// TODO: Use the new Command attribute of _ModUserInfo

type _MGI _ModGenInfo
var (
	realMain Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGI]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGI](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGI])

		var p_user_location *ULComm.UserLocation = ULComm.GetUserLocation()
		var user_location ULComm.UserLocation
		if p_user_location != nil {
			user_location = *p_user_location
		}

		var notifs_were_true map[string]bool = make(map[string]bool)

		var last_md5 []byte = nil
		var prev_curr_last_known_user_loc string = user_location.Curr_location
		var prev_prev_last_known_user_loc string = user_location.Prev_location
		for {
			var new_md5 []byte = Utils.GetFileContentsWEBSITE("reminders.json", true)
			if new_md5 != nil && !bytes.Equal(new_md5, last_md5) {
				updateLocalReminders()

				last_md5 = new_md5
			}

			var reminders []RRComm.Reminder = moduleInfo_GL.ModGenInfo.Reminders

			// Add each reminder to the internal reminders list
			var list_modified bool = false
			var reminders_info_list map[string]int64 = moduleInfo_GL.ModGenInfo.Reminders_info
			if reminders_info_list == nil {
				reminders_info_list = make(map[string]int64)
				moduleInfo_GL.ModGenInfo.Reminders_info = reminders_info_list
				list_modified = true
			}
			for _, reminder := range reminders {
				if _, ok := reminders_info_list[reminder.Id]; !ok {
					reminders_info_list[reminder.Id] = 0
					list_modified = true
				}
			}
			if list_modified {
				_ = moduleInfo_GL.UpdateGenInfo()
			}

			// Location trigger - if the user location changed, check if any reminder is triggered
			p_user_location = ULComm.GetUserLocation()
			if p_user_location != nil {
				user_location = *p_user_location
			}
			var curr_last_known_user_loc string = user_location.Curr_location
			var prev_last_known_user_loc string = user_location.Prev_location
			if curr_last_known_user_loc != prev_curr_last_known_user_loc || prev_last_known_user_loc != prev_prev_last_known_user_loc {
				prev_curr_last_known_user_loc = curr_last_known_user_loc
				prev_prev_last_known_user_loc = prev_last_known_user_loc

				for _, reminder := range reminders {
					// If the reminder has a time set or has no location, skip it
					if reminder.Time != "" || reminder.User_location == "" {
						continue
					}

					// In case there's a "+", the user must have arrived at the location. If there's a "-", the user
					// must have left the location.
					var condition_loc bool
					if strings.HasPrefix(reminder.User_location, "+") {
						var rem_loc string = reminder.User_location[1:]
						condition_loc = checkLocation(rem_loc, curr_last_known_user_loc)
					} else if strings.HasPrefix(reminder.User_location, "-") {
						var rem_loc string = reminder.User_location[1:]
						condition_loc = checkLocation(rem_loc, prev_last_known_user_loc)
					} else {
						// Nothing to do
						continue
					}

					var condition bool = checkCondition(reminder, notifs_were_true)

					if condition_loc && condition {
						MOD_3.QueueSpeech(reminder.Message, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE1_ALWAYS_NOTIFY)

						log.Println("Reminder! Message: " + reminder.Message)
					}
				}
			}

			// Time/condition trigger - if the time changed (it always does), check if any reminder is triggered
			for _, reminder := range reminders {
				var condition_time bool = false
				var test_time int64 = 0
				// If the reminder has no time set, skip it
				if reminder.Time != "" {
					var curr_time int64 = time.Now().Unix() / 60
					var reminder_time string = reminder.Time
					var format string = "2006-01-02 -- 15:04:05"
					t, _ := time.ParseInLocation(format, reminder_time, time.Local)
					test_time = t.Unix() / 60
					if reminder.Repeat_each > 0 {
						var repeat_each int64 = reminder.Repeat_each
						for {
							if test_time + repeat_each <= curr_time {
								test_time += repeat_each
							} else {
								break
							}
						}
					}

					condition_time  = curr_time >= test_time && reminders_info_list[reminder.Id] != test_time
				} else {
					condition_time = true
				}

				// Check if the reminder is due and if it was already reminded

				var condition_loc bool = false
				if reminder.User_location != "" {
					// Check if the reminder has a location and the user is at that location.
					var curr_user_loc string = user_location.Curr_location
					if curr_user_loc != MOD_12.UNKNOWN_LOCATION {
						condition_loc = checkLocation(reminder.User_location, curr_user_loc)
					}
				} else {
					condition_loc = true
				}

				var condition bool = checkCondition(reminder, notifs_were_true)

				if condition_time && condition_loc && condition {
					MOD_3.QueueSpeech(reminder.Message, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE1_ALWAYS_NOTIFY)

					log.Println("Reminder! Message: " + reminder.Message)

					// Set the last reminded time to the test time
					reminders_info_list[reminder.Id] = test_time

					_ = moduleInfo_GL.UpdateGenInfo()
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, TIME_SLEEP_S) {
				return
			}
		}
	}
}

func checkLocation(reminder_loc string, location string) bool {
	if strings.HasSuffix(reminder_loc, "*") {
		// If the reminder location ends with a "*", it means that the user must be at a location that starts with the
		// reminder location.
		reminder_loc = reminder_loc[:len(reminder_loc) - 1]

		return strings.HasPrefix(location, reminder_loc)
	}

	return reminder_loc == location
}

func computeCondition(condition string) bool {
	condition = formatCondition(condition)
	//log.Println("Condition:", condition)
	expr, err := eval.ParseString(condition, "")
	if err != nil {
		return false
	}
	r, err := expr.EvalToInterface(nil)
	if err != nil {
		return false
	}

	return r.(bool)
}

func formatCondition(condition string) string {
	var power_connected bool = Registry.GetValue(ClientRegKeys.K_POWER_CONNECTED).GetData(true, nil).(bool)
	var battery_level int = Registry.GetValue(ClientRegKeys.K_BATTERY_LEVEL).GetData(true, nil).(int)
	var screen_brightness int = Registry.GetValue(ClientRegKeys.K_SCREEN_BRIGHTNESS).GetData(true, nil).(int)
	var sound_volume int = Registry.GetValue(ClientRegKeys.K_SOUND_VOLUME).GetData(true, nil).(int)
	var sound_muted bool = Registry.GetValue(ClientRegKeys.K_SOUND_MUTED).GetData(true, nil).(bool)

	condition = strings.Replace(condition, "power_connected", strconv.FormatBool(power_connected), -1)
	condition = strings.Replace(condition, "battery_level", strconv.Itoa(battery_level), -1)
	condition = strings.Replace(condition, "screen_brightness", strconv.Itoa(screen_brightness), -1)
	condition = strings.Replace(condition, "sound_volume", strconv.Itoa(sound_volume), -1)
	condition = strings.Replace(condition, "sound_muted", strconv.FormatBool(sound_muted), -1)

	return condition
}

func checkCondition(reminder RRComm.Reminder, notifs_were_true map[string]bool) bool {
	var condition bool = false
	if reminder.Device_condition != "" {
		if ok := notifs_were_true[reminder.Id]; !ok {
			notifs_were_true[reminder.Id] = false
		}

		if computeCondition(reminder.Device_condition) {
			if !notifs_were_true[reminder.Id] {
				notifs_were_true[reminder.Id] = true

				condition = true
			}
		} else {
			notifs_were_true[reminder.Id] = false
		}
	} else {
		condition = true
	}

	return condition
}

func updateLocalReminders() {
	var p_reminders *[]RRComm.Reminder = RRComm.GetRemindersList()
	if p_reminders == nil {
		return
	}

	moduleInfo_GL.ModGenInfo.Reminders = *p_reminders
	_ = moduleInfo_GL.UpdateGenInfo()
}
