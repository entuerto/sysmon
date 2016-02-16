// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

const (
	PDH_CSTATUS_VALID_DATA                      = 0x00000000 
	PDH_CSTATUS_NEW_DATA                        = 0x00000001 
	PDH_CSTATUS_NO_MACHINE                      = 0x800007D0 
	PDH_CSTATUS_NO_INSTANCE                     = 0x800007D1 
	PDH_MORE_DATA                               = 0x800007D2 
	PDH_CSTATUS_ITEM_NOT_VALIDATED              = 0x800007D3 
	PDH_RETRY                                   = 0x800007D4 
	PDH_NO_DATA                                 = 0x800007D5 
	PDH_CALC_NEGATIVE_DENOMINATOR               = 0x800007D6 
	PDH_CALC_NEGATIVE_TIMEBASE                  = 0x800007D7 
	PDH_CALC_NEGATIVE_VALUE                     = 0x800007D8 
	PDH_DIALOG_CANCELLED                        = 0x800007D9 
	PDH_END_OF_LOG_FILE                         = 0x800007DA 
	PDH_ASYNC_QUERY_TIMEOUT                     = 0x800007DB 
	PDH_CANNOT_SET_DEFAULT_REALTIME_DATASOURCE  = 0x800007DC 
	PDH_CSTATUS_NO_OBJECT                       = 0xC0000BB8 
	PDH_CSTATUS_NO_COUNTER                      = 0xC0000BB9 
	PDH_CSTATUS_INVALID_DATA                    = 0xC0000BBA 
	PDH_MEMORY_ALLOCATION_FAILURE               = 0xC0000BBB 
	PDH_INVALID_HANDLE                          = 0xC0000BBC 
	PDH_INVALID_ARGUMENT                        = 0xC0000BBD 
	PDH_FUNCTION_NOT_FOUND                      = 0xC0000BBE 
	PDH_CSTATUS_NO_COUNTERNAME                  = 0xC0000BBF 
	PDH_CSTATUS_BAD_COUNTERNAME                 = 0xC0000BC0 
	PDH_INVALID_BUFFER                          = 0xC0000BC1 
	PDH_INSUFFICIENT_BUFFER                     = 0xC0000BC2 
	PDH_CANNOT_CONNECT_MACHINE                  = 0xC0000BC3 
	PDH_INVALID_PATH                            = 0xC0000BC4 
	PDH_INVALID_INSTANCE                        = 0xC0000BC5 
	PDH_INVALID_DATA                            = 0xC0000BC6 
	PDH_NO_DIALOG_DATA                          = 0xC0000BC7 
	PDH_CANNOT_READ_NAME_STRINGS                = 0xC0000BC8 
	PDH_LOG_FILE_CREATE_ERROR                   = 0xC0000BC9 
	PDH_LOG_FILE_OPEN_ERROR                     = 0xC0000BCA 
	PDH_LOG_TYPE_NOT_FOUND                      = 0xC0000BCB 
	PDH_NO_MORE_DATA                            = 0xC0000BCC 
	PDH_ENTRY_NOT_IN_LOG_FILE                   = 0xC0000BCD 
	PDH_DATA_SOURCE_IS_LOG_FILE                 = 0xC0000BCE 
	PDH_DATA_SOURCE_IS_REAL_TIME                = 0xC0000BCF 
	PDH_UNABLE_READ_LOG_HEADER                  = 0xC0000BD0 
	PDH_FILE_NOT_FOUND                          = 0xC0000BD1 
	PDH_FILE_ALREADY_EXISTS                     = 0xC0000BD2 
	PDH_NOT_IMPLEMENTED                         = 0xC0000BD3 
	PDH_STRING_NOT_FOUND                        = 0xC0000BD4 
	PDH_UNABLE_MAP_NAME_FILES                   = 0x80000BD5 
	PDH_UNKNOWN_LOG_FORMAT                      = 0xC0000BD6 
	PDH_UNKNOWN_LOGSVC_COMMAND                  = 0xC0000BD7 
	PDH_LOGSVC_QUERY_NOT_FOUND                  = 0xC0000BD8 
	PDH_LOGSVC_NOT_OPENED                       = 0xC0000BD9 
	PDH_WBEM_ERROR                              = 0xC0000BDA 
	PDH_ACCESS_DENIED                           = 0xC0000BDB 
	PDH_LOG_FILE_TOO_SMALL                      = 0xC0000BDC 
	PDH_INVALID_DATASOURCE                      = 0xC0000BDD 
	PDH_INVALID_SQLDB                           = 0xC0000BDE 
	PDH_NO_COUNTERS                             = 0xC0000BDF 
	PDH_SQL_ALLOC_FAILED                        = 0xC0000BE0 
	PDH_SQL_ALLOCCON_FAILED                     = 0xC0000BE1 
	PDH_SQL_EXEC_DIRECT_FAILED                  = 0xC0000BE2 
	PDH_SQL_FETCH_FAILED                        = 0xC0000BE3 
	PDH_SQL_ROWCOUNT_FAILED                     = 0xC0000BE4 
	PDH_SQL_MORE_RESULTS_FAILED                 = 0xC0000BE5 
	PDH_SQL_CONNECT_FAILED                      = 0xC0000BE6 
	PDH_SQL_BIND_FAILED                         = 0xC0000BE7 
	PDH_CANNOT_CONNECT_WMI_SERVER               = 0xC0000BE8 
	PDH_PLA_COLLECTION_ALREADY_RUNNING          = 0xC0000BE9 
	PDH_PLA_ERROR_SCHEDULE_OVERLAP              = 0xC0000BEA 
	PDH_PLA_COLLECTION_NOT_FOUND                = 0xC0000BEB 
	PDH_PLA_ERROR_SCHEDULE_ELAPSED              = 0xC0000BEC 
	PDH_PLA_ERROR_NOSTART                       = 0xC0000BED 
	PDH_PLA_ERROR_ALREADY_EXISTS                = 0xC0000BEE 
	PDH_PLA_ERROR_TYPE_MISMATCH                 = 0xC0000BEF 
	PDH_PLA_ERROR_FILEPATH                      = 0xC0000BF0 
	PDH_PLA_SERVICE_ERROR                       = 0xC0000BF1 
	PDH_PLA_VALIDATION_ERROR                    = 0xC0000BF2 
	PDH_PLA_VALIDATION_WARNING                  = 0x80000BF3 
	PDH_PLA_ERROR_NAME_TOO_LONG                 = 0xC0000BF4 
	PDH_INVALID_SQL_LOG_FORMAT                  = 0xC0000BF5 
	PDH_COUNTER_ALREADY_IN_QUERY                = 0xC0000BF6 
	PDH_BINARY_LOG_CORRUPT                      = 0xC0000BF7 
	PDH_LOG_SAMPLE_TOO_SMALL                    = 0xC0000BF8 
	PDH_OS_LATER_VERSION                        = 0xC0000BF9 
	PDH_OS_EARLIER_VERSION                      = 0xC0000BFA 
	PDH_INCORRECT_APPEND_TIME                   = 0xC0000BFB 
	PDH_UNMATCHED_APPEND_COUNTER                = 0xC0000BFC 
	PDH_SQL_ALTER_DETAIL_FAILED                 = 0xC0000BFD 
	PDH_QUERY_PERF_DATA_TIMEOUT                 = 0xC0000BFE 
)

var codeText = map[uintptr]string{
	PDH_CSTATUS_VALID_DATA                      : "The returned data is valid.", 
	PDH_CSTATUS_NEW_DATA                        : "The return data value is valid and different from the last sample.",
	PDH_CSTATUS_NO_MACHINE                      : "Unable to connect to the specified computer, or the computer is offline.",
	PDH_CSTATUS_NO_INSTANCE                     : "The specified instance is not present.",
	PDH_MORE_DATA                               : "There is more data to return than would fit in the supplied buffer. Allocate a larger buffer and call the function again.",
	PDH_CSTATUS_ITEM_NOT_VALIDATED              : "The data item has been added to the query but has not been validated nor accessed. No other status information on this data item is available.",
	PDH_RETRY                                   : "The selected operation should be retried.",
	PDH_NO_DATA                                 : "No data to return.",
	PDH_CALC_NEGATIVE_DENOMINATOR               : "A counter with a negative denominator value was detected.",
	PDH_CALC_NEGATIVE_TIMEBASE                  : "A counter with a negative time base value was detected.",
	PDH_CALC_NEGATIVE_VALUE                     : "A counter with a negative value was detected.",
	PDH_DIALOG_CANCELLED                        : "The user canceled the dialog box.",
	PDH_END_OF_LOG_FILE                         : "The end of the log file was reached.",
	PDH_ASYNC_QUERY_TIMEOUT                     : "A time-out occurred while waiting for the asynchronous counter collection thread to end.",
	PDH_CANNOT_SET_DEFAULT_REALTIME_DATASOURCE  : "Cannot change set default real-time data source. There are real-time query sessions collecting counter data.",
	PDH_CSTATUS_NO_OBJECT                       : "The specified object is not found on the system.",
	PDH_CSTATUS_NO_COUNTER                      : "The specified counter could not be found.",
	PDH_CSTATUS_INVALID_DATA                    : "The returned data is not valid.",
	PDH_MEMORY_ALLOCATION_FAILURE               : "A PDH function could not allocate enough temporary memory to complete the operation. Close some applications or extend the page file and retry the function.",
	PDH_INVALID_HANDLE                          : "The handle is not a valid PDH object.",
	PDH_INVALID_ARGUMENT                        : "A required argument is missing or incorrect.",
	PDH_FUNCTION_NOT_FOUND                      : "Unable to find the specified function.",
	PDH_CSTATUS_NO_COUNTERNAME                  : "No counter was specified.", 
	PDH_CSTATUS_BAD_COUNTERNAME                 : "Unable to parse the counter path. Check the format and syntax of the specified path.",
	PDH_INVALID_BUFFER                          : "The buffer passed by the caller is not valid.",
	PDH_INSUFFICIENT_BUFFER                     : "The requested data is larger than the buffer supplied. Unable to return the requested data.",
	PDH_CANNOT_CONNECT_MACHINE                  : "Unable to connect to the requested computer.",
	PDH_INVALID_PATH                            : "The specified counter path could not be interpreted.",
	PDH_INVALID_INSTANCE                        : "The instance name could not be read from the specified counter path.",
	PDH_INVALID_DATA                            : "The data is not valid.",
	PDH_NO_DIALOG_DATA                          : "The dialog box data block was missing or not valid.",
	PDH_CANNOT_READ_NAME_STRINGS                : "Unable to read the counter and/or help text from the specified computer.",
	PDH_LOG_FILE_CREATE_ERROR                   : "Unable to create the specified log file.",
	PDH_LOG_FILE_OPEN_ERROR                     : "Unable to open the specified log file.",
	PDH_LOG_TYPE_NOT_FOUND                      : "The specified log file type has not been installed on this system.",
	PDH_NO_MORE_DATA                            : "No more data is available.",
	PDH_ENTRY_NOT_IN_LOG_FILE                   : "The specified record was not found in the log file.",
	PDH_DATA_SOURCE_IS_LOG_FILE                 : "The specified data source is a log file.",
	PDH_DATA_SOURCE_IS_REAL_TIME                : "The specified data source is the current activity.",
	PDH_UNABLE_READ_LOG_HEADER                  : "The log file header could not be read.",
	PDH_FILE_NOT_FOUND                          : "Unable to find the specified file.",
	PDH_FILE_ALREADY_EXISTS                     : "There is already a file with the specified file name.",
	PDH_NOT_IMPLEMENTED                         : "The function referenced has not been implemented.",
	PDH_STRING_NOT_FOUND                        : "Unable to find the specified string in the list of performance name and help text strings.",
	PDH_UNABLE_MAP_NAME_FILES                   : "Unable to map to the performance counter name data files. The data will be read from the registry and stored locally.",
	PDH_UNKNOWN_LOG_FORMAT                      : "The format of the specified log file is not recognized by the PDH DLL.",
	PDH_UNKNOWN_LOGSVC_COMMAND                  : "The specified Log Service command value is not recognized.",
	PDH_LOGSVC_QUERY_NOT_FOUND                  : "The specified query from the Log Service could not be found or could not be opened.",
	PDH_LOGSVC_NOT_OPENED                       : "The Performance Data Log Service key could not be opened. This may be due to insufficient privilege or because the service has not been installed.",
	PDH_WBEM_ERROR                              : "An error occurred while accessing the WBEM data store.",
	PDH_ACCESS_DENIED                           : "Unable to access the desired computer or service. Check the permissions and authentication of the log service or the interactive user session against those on the computer or service being monitored.",
	PDH_LOG_FILE_TOO_SMALL                      : "The maximum log file size specified is too small to log the selected counters. No data will be recorded in this log file. Specify a smaller set of counters to log or a larger file size and retry this call.",
	PDH_INVALID_DATASOURCE                      : "Cannot connect to ODBC DataSource Name.", 
	PDH_INVALID_SQLDB                           : "SQL Database does not contain a valid set of tables for Perfmon.",
	PDH_NO_COUNTERS                             : "No counters were found for this Perfmon SQL Log Set.",
	PDH_SQL_ALLOC_FAILED                        : "Call to SQLAllocStmt failed with %1.",
	PDH_SQL_ALLOCCON_FAILED                     : "Call to SQLAllocConnect failed with %1.",
	PDH_SQL_EXEC_DIRECT_FAILED                  : "Call to SQLExecDirect failed with %1.",
	PDH_SQL_FETCH_FAILED                        : "Call to SQLFetch failed with %1.",
	PDH_SQL_ROWCOUNT_FAILED                     : "Call to SQLRowCount failed with %1.",
	PDH_SQL_MORE_RESULTS_FAILED                 : "Call to SQLMoreResults failed with %1.",
	PDH_SQL_CONNECT_FAILED                      : "Call to SQLConnect failed with %1.",
	PDH_SQL_BIND_FAILED                         : "Call to SQLBindCol failed with %1.",
	PDH_CANNOT_CONNECT_WMI_SERVER               : "Unable to connect to the WMI server on requested computer.",
	PDH_PLA_COLLECTION_ALREADY_RUNNING          : "Collection \"%1!s!\" is already running.",
	PDH_PLA_ERROR_SCHEDULE_OVERLAP              : "The specified start time is after the end time.",
	PDH_PLA_COLLECTION_NOT_FOUND                : "Collection \"%1!s!\" does not exist.",
	PDH_PLA_ERROR_SCHEDULE_ELAPSED              : "The specified end time has already elapsed.",
	PDH_PLA_ERROR_NOSTART                       : "Collection \"%1!s!\" did not start; check the application event log for any errors.",
	PDH_PLA_ERROR_ALREADY_EXISTS                : "Collection \"%1!s!\" already exists.",
	PDH_PLA_ERROR_TYPE_MISMATCH                 : "There is a mismatch in the settings type.",
	PDH_PLA_ERROR_FILEPATH                      : "The information specified does not resolve to a valid path name.",
	PDH_PLA_SERVICE_ERROR                       : "The \"Performance Logs & Alerts\" service did not respond.",
	PDH_PLA_VALIDATION_ERROR                    : "The information passed is not valid.",
	PDH_PLA_VALIDATION_WARNING                  : "The information passed is not valid.",
	PDH_PLA_ERROR_NAME_TOO_LONG                 : "The name supplied is too long.",
	PDH_INVALID_SQL_LOG_FORMAT                  : "SQL log format is incorrect. Correct format is \"SQL:<DSN-name>!<LogSet-Name>\".",
	PDH_COUNTER_ALREADY_IN_QUERY                : "Performance counter in PdhAddCounter call has already been added in the performance query. This counter is ignored.",
	PDH_BINARY_LOG_CORRUPT                      : "Unable to read counter information and data from input binary log files.",
	PDH_LOG_SAMPLE_TOO_SMALL                    : "At least one of the input binary log files contain fewer than two data samples.",
	PDH_OS_LATER_VERSION                        : "The version of the operating system on the computer named %1 is later than that on the local computer. This operation is not available from the local computer.",
	PDH_OS_EARLIER_VERSION                      : "%1 supports %2 or later. Check the operating system version on the computer named %3.",
	PDH_INCORRECT_APPEND_TIME                   : "The output file must contain earlier data than the file to be appended.",
	PDH_UNMATCHED_APPEND_COUNTER                : "Both files must have identical counters in order to append.",
	PDH_SQL_ALTER_DETAIL_FAILED                 : "Cannot alter CounterDetail table layout in SQL database.",
	PDH_QUERY_PERF_DATA_TIMEOUT                 : "System is busy. A time-out occurred when collecting counter data. Please retry later or increase the CollectTime registry value.",
}