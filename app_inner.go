//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (b *App) errorChanListen(ctx context.Context) {
	for {
		select {
		case err := <-b.errorChannel:
			panicWithErrorMessageDialog(ctx, err)
		case <-b.stopChannel:
			return
		}
	}
}

func (b *App) stopChanListen(ctx context.Context) {
	for {
		select {
		case <-b.stopChannel:
			return
		}
	}
}

func panicWithErrorMessageDialog(ctx context.Context, err error) {
	_, dialogErr := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Title:         "Error",
		Type:          runtime.ErrorDialog,
		Message:       err.Error(),
		Buttons:       []string{"OK"},
		DefaultButton: "OK",
	})
	if dialogErr != nil {
		fmt.Printf("[CRITIC] open dialog error %s\n", dialogErr.Error())
	}
}
