package repository_errors

import "errors"

var ErrRoomAlreadyExists = errors.New("room with this name already exists")
var ErrRoomScheduleExists = errors.New("schedule for this room is already exists")
var ErrRoomNotFound = errors.New("room with this id wasn't found")
var ErrBookingNoExists = errors.New("booking with such id is not exist")
var ErrSlotIdConflict = errors.New("this slot is already busy")
var ErrSlotNotFound = errors.New("slot with such id is not exist")
var ErrUserNotFound = errors.New("user with such email not found")
