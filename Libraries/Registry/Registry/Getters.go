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

package Registry

import (
	"strconv"
)

/*
getInternal returns whether the data is internal and the data to return.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data
  - no_data – the data to return if there's no data or nil to return the default values

– Returns:
  - whether the no_data parameter was used
  - the data to return
 */
func (value *Value) getInternal(curr_data bool, no_data any) (bool, any) {
	if no_data != nil {
		if curr_data {
			if value.time_updated_curr == 0 {
				return true, no_data
			}
		} else {
			if value.time_updated_prev == 0 {
				return true, no_data
			}
		}
	}

	return false, nil
}

/*
GetTimeUpdated returns the time the data was updated in milliseconds.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the time the data was updated in milliseconds
 */
func (value *Value) GetTimeUpdated(curr_data bool) int64 {
	if curr_data {
		return value.time_updated_curr
	} else {
		return value.time_updated_prev
	}
}

/*
GetType returns the type of the Value.

-----------------------------------------------------------

– Returns:
  - the type of the Value
 */
func (value *Value) GetType() string {
	return value.type_
}

/*
GetBool returns the boolean value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the boolean value of the Value
 */
func (value *Value) GetBool(curr_data bool) bool {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseBool(data)
	if err != nil {
		return false
	}

	return i
}

/*
GetInt returns the integer value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the integer value of the Value
 */
func (value *Value) GetInt(curr_data bool) int {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.Atoi(data)
	if err != nil {
		return -1
	}

	return i
}

/*
GetLong returns the long value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the long value of the Value
 */
func (value *Value) GetLong(curr_data bool) int64 {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return -1
	}

	return i
}

/*
GetFloat returns the float value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the float value of the Value
 */
func (value *Value) GetFloat(curr_data bool) float32 {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseFloat(data, 32)
	if err != nil {
		return -1
	}

	return float32(i)
}

/*
GetDouble returns the double value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the double value of the Value
 */
func (value *Value) GetDouble(curr_data bool) float64 {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return -1
	}

	return i
}

/*
GetString returns the string value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the string value of the Value
 */
func (value *Value) GetString(curr_data bool) string {
	if curr_data {
		return value.curr_data
	} else {
		return value.prev_data
	}
}

/*
GetData returns the data of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data
  - no_data – the data to return if there's no data or nil to return the default values
 */
func (value *Value) GetData(curr_data bool, no_data any) any {
	no_data_used, no_data_ret := value.getInternal(curr_data, no_data)
	if no_data_used {
		return no_data_ret
	}

	switch value.type_ {
		case TYPE_BOOL:
			return value.GetBool(curr_data)
		case TYPE_INT:
			return value.GetInt(curr_data)
		case TYPE_LONG:
			return value.GetLong(curr_data)
		case TYPE_FLOAT:
			return value.GetFloat(curr_data)
		case TYPE_DOUBLE:
			return value.GetDouble(curr_data)
		case TYPE_STRING:
			return value.GetString(curr_data)
	}

	// Won't happen
	return nil
}
