# NAME
**esp_app_trace_viewer** - Viewer for ESP32 application tracing

# SYNOPSIS
**esp_app_trace_viewer**
[*interface*]
[*target*]
[*filename*]

# DESCRIPTION
**esp_app_trace_viewer**
is a utility for viewing ESP32 application tracing using OpenOCD for the ESP32.

The example batch and bash script files set up the viewer for ESP32S3 using an ftdi JTAG adapter.

When viewing 'r' enter will reset target, 'x' enter will exit viewer.

# EXAMPLES

View application trace from ESP32S3 using ftdi JTAG adapter, trace data will be written to trace.txt:

	esp_app_trace_viewer interface/ftdi/esp32_devkitj_v1.cfg target/esp32s3.cfg trace.txt

# AUTHORS
sjp27 &lt; https://github.com/sjp27 &gt;
implemented esp_app_trace_viewer.

ESP32 is a trademark of Espressif Systems (Shanghai) Co., Ltd.
