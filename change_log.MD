# Schedule_Task_Command_Time  version 2.15.1

### General
- use config to control swagger version and host
- use logrus as new log 
- modify README.MD


# Schedule_Task_Command_Time  version 2.14.0

### General
- config file update
- util config function update
- golang and package update

### Logs
- read api update
- write api add parameters
  - account
  - ip
  - method
  - module
  - status_code

### all server
- add locals variables
------------------------------------------------------------------


# Schedule_Task_Command_Time  version 2.12.3

### redis add cluster connection

### Websocket hub
- improve websocket
------------------------------------------------------------------


# Schedule_Task_Command_Time  version 2.11.0

### Websocket hub
- add command websocket
- improve websocket hub manager
------------------------------------------------------------------


# Schedule_Task_Command_Time  version 2.10.0

### Command server
- add default variable to command
  - key: __command_id__
- add default parser return
  - key: __status_code__
  - key: __command_id__
  - key: __Status__

### Task server
- can send global variables to task(set int as -1)
  - ex: variables:{"-1":{"global_event_id": 5487}}
- add default global variables
  - key: __task_id__
------------------------------------------------------------------


# Schedule_Task_Command_Time  version 2.9.1

### all API
- add log write to influxDB

### log API
- add get historical logs
------------------------------------------------------------------


# Schedule_Task_Command_Time  version 2.8.0

### Task Server
- Get history data add new tag "id"

### Task API
- modify history data, add query id
- get task by id can search both live data and history data
- get simple task by id can search both live data and history data
- get stage_item status can search both live data and history data

### Command Server
- Get history data add new tag "id"

### Command API
- modify history data, add query id
- get command by id can search both live data and history data

### Time Server
- Get history data add new tag "id"

### Time API
- modify history data, add query id
------------------------------------------------------------------


# Schedule_Task_Command_Time  version 2.7.17

### Task Server
- command parser to variables

### Task API
- add simple task api
- add simple stages status api

### Command Server
- parser return data
------------------------------------------------------------------
## 下一個版本預計完成


# Schedule_Task_Command_Time  version 2.6.0

### Task Server
- redis stream read
- redis stream add

### Command Server
- redis stream read
- redis stream add
------------------------------------------------------------------
## 下一個版本預計完成
- command parser
- task change variable


# Schedule_Task_Command_Time  version 2.5.7

### Task Server
- Cancel task 增加 message
- history 可搜尋全部狀態
- stages加上變數(Variable)
- 刪除taskId
- 新增ID
- task stages中新增stage_id

### Command Server
- Cancel command 增加 message
- history 可搜尋全部狀態
- 刪除commandId
- 新增ID

### Task Template API
- mCondition id 顯示

### Time API
- [GET] history

### Time Server
- redis topic 修改傳入格式
------------------------------------------------------------------
## 下一個版本預計完成



# Schedule_Task_Command_Time  version 2.4.0

### Task Server
- 新增變數

### Command Server
- 新增變數
- DB command_template新增variable欄位
- command_template add Patch method 修改command API
------------------------------------------------------------------
## 下一個版本預計完成
 - cancel task 新增message
 - cancel command 新增message
 - Gamma 文件


# Schedule_Task_Command_Time  version 2.3.0

### Time server
- 觸發Time template的time data
- Redis topic觸發功能

### Schedule Server
- 觸發Schedule

### Task Server
- 執行Task api
- 歷史Task api
- 即時Task api
- Redis topic觸發功能

### Command Server
- 執行Command api
- 歷史Command api
- 即時Command api
- Redis topic觸發功能
------------------------------------------------------------------
## 下一個版本預計完成
Task 增加變數
Command 增加變數