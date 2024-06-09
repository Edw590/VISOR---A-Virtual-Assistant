/*******************************************************************************
 * Copyright 2023-2023 Edw590
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

// Package main provides various examples of Fyne API capabilities.
package main

import (
	"GPT/GPT"
	"OIG/OIG"
	"Utils"
	"VISOR/AppCode/Screens"
	"VISOR/logo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"strings"
	"time"
)

var my_app_GL fyne.App = nil

func main() {
	Utils.PersonalConsts_GL.Init()
	GPT.SetWebsiteInfo(Utils.PersonalConsts_GL.WEBSITE_URL, Utils.PersonalConsts_GL.WEBSITE_PW)
	OIG.SetWebsiteInfo(Utils.PersonalConsts_GL.WEBSITE_URL, Utils.PersonalConsts_GL.WEBSITE_PW)

	// Create a new application
	my_app_GL = app.NewWithID("com.edw590.visor_c")
	my_app_GL.SetIcon(logo.LogoBlackGmail)
	var my_window fyne.Window = my_app_GL.NewWindow("V.I.S.O.R.")

	processNotifications()


	// Create the content area with a label to display different screens
	var content_label *widget.Label = widget.NewLabel("Welcome!")
	var content_container *fyne.Container = container.NewVBox(content_label)

	// Create the navigation bar
	var nav_bar *fyne.Container = container.NewVBox(
		widget.NewButton("Home", func() {
			content_container.Objects = []fyne.CanvasObject{Screens.Home()}
			content_container.Refresh()
		}),
		widget.NewButton("Dev Mode", func() {
			content_container.Objects = []fyne.CanvasObject{Screens.DevMode()}
			content_container.Refresh()
		}),
		widget.NewButton("Communicator", func() {
			content_container.Objects = []fyne.CanvasObject{Screens.Communicator()}
			content_container.Refresh()
		}),
		/*widget.NewButton("Progress Bar", func() {
			contentContainer.Objects = []fyne.CanvasObject{createTextScreen()}
			contentContainer.Refresh()
		}),*/
	)


	// Create a split container to hold the navigation bar and the content
	var split *container.Split = container.NewHSplit(nav_bar, content_container)
	split.SetOffset(0.2) // Set the split ratio (20% for nav, 80% for content)

	// Set the content of the window
	my_window.SetContent(split)

	var prev_screen fyne.CanvasObject = Screens.Current_screen_GL
	// Add system tray functionality
	if desk, ok := my_app_GL.(desktop.App); ok {
		var icon *fyne.StaticResource = logo.LogoBlackGmail
		var menu *fyne.Menu = fyne.NewMenu("Tray",
			fyne.NewMenuItem("Show", func() {
				// Hide too because in case the window is shown but behind other apps, it won't show. So hiding and
				// showing does it. Maybe this happens because RequestFocus doesn't always work? Who knows. But this
				// fixes whatever the problem is.
				my_window.Hide()
				my_window.Show()
				my_window.RequestFocus()

				// Restore the previous screen state
				Screens.Current_screen_GL = prev_screen
			}),
		)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(icon)
	}

	// Minimize to tray on close
	my_window.SetCloseIntercept(func() {
		// Store the previous screen before hiding
		Screens.Current_screen_GL = nil
		my_window.Hide()
	})

	// Show and run the application
	my_window.Resize(fyne.NewSize(640, 480))
	my_window.ShowAndRun()
}

/*
processNotifications processes in a different thread the notifications queued in the notifications folder.
 */
func processNotifications() {
	go func() {
		for {
			var file_list []Utils.FileInfo = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_VISOR).
				Add2(true, Utils.NOTIFS_REL_FOLDER).GetFileList()
			for _, file := range file_list {
				// Display the notification
				notification := fyne.NewNotification(strings.Split(file.Name, "-")[0], *file.GPath.ReadTextFile())
				my_app_GL.SendNotification(notification)

				// Remove the file
				_ = file.GPath.Remove()

				time.Sleep(5 * time.Second)
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
