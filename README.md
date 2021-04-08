# Go-Pichubot

> 皮丘Bot in Golang

![](https://img.shields.io/static/v1?label=Pichubot&message=Golang&color=?style=flat&logo=data:image/png;base64,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c6XHXENU4Fh/Krbvi2iPv1yq8LlkR121RHyVBAW2gPf+E+3+7Za+Dcvnlx589uWDV+pAmduGWFL1N/strZJIgGtIKssQJLKaaUFj8S/uWuGq9pXzZHhQR/vnY39+7b0vXuZ5zl/rLWppqCyaO7fJ6jfaSHTz9oFHFt/o+PSl149eLURDwShWLqlZOGuaqXSiaz3PES4ez+pCnmVOn1bQ3ZP74p19Q6+67JagoqrBk33hvSf7wpbte+EGEVIb75wSeOVv3VcrIcMwIITYn314elu5i7vm14/XlKLCZXDzWUHbvTtc2N49vJmw2P77H887t/5XO746FRMCzyVkmcOrb3ZfEeJSE/qnOOF22ZpnN9qvF4QJjlNK0VBrNHrsgtzdF08NDEXfIjr23fPrnVeNZEUtIBrLXBXiEoDbV8yCxSRW2S2cmao6qE7G16+LBwh0nUDXx7v/4sgtKmKQL9EO7zoVvLegqcdPng1NDPzXS/Dm2wfRPxQN7j8Yi7vgKjIaCQwyAWGAgg6oeQI1r4MhDCilYHnAaCFgJYqaCrl/LJN612wQwDATX0YuAWic5sNYLOO8/Zb679kckpRSKeIBHb2hXDaSUPqGY4qSTKvGYpPgrPPKZq/DwFplFpk4g4iSR0dnNPkgIdINzVXZ2iontr57aHIAHtmwCL95edfSslJ5NUfI2V0dkei54cxgV09sd8exob0AVQFYJMlQ4bEbayuc8oypVbYlrY328poqQbjj1rI7rTbZ9mHH8NPh0eSZSVFw0ViWccuy1ATwNQDsAIQXH10MhuEuvGfhdljw4sNLAEACyDyfx/rIsvnlbzy5YfrJ555ojcyeWfnA4xsXTyrvfwHMjdA8JZE39QAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMC0wNy0xOVQwMzozOToyMCswMDowMIZnDlwAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMTktMDEtMDhUMTc6NDI6NTYrMDA6MDB+3sVrAAAAIHRFWHRzb2Z0d2FyZQBodHRwczovL2ltYWdlbWFnaWNrLm9yZ7zPHZ0AAAAYdEVYdFRodW1iOjpEb2N1bWVudDo6UGFnZXMAMaf/uy8AAAAYdEVYdFRodW1iOjpJbWFnZTo6SGVpZ2h0ADI1NunDRBkAAAAXdEVYdFRodW1iOjpJbWFnZTo6V2lkdGgAMjU2ejIURAAAABl0RVh0VGh1bWI6Ok1pbWV0eXBlAGltYWdlL3BuZz+yVk4AAAAXdEVYdFRodW1iOjpNVGltZQAxNTQ2OTY5Mzc2rTseLwAAABJ0RVh0VGh1bWI6OlNpemUANzgxMTlCQrbM/gAAAFh0RVh0VGh1bWI6OlVSSQBmaWxlOi8vL2RhdGEvd3d3cm9vdC93d3cuZWFzeWljb24ubmV0L2Nkbi1pbWcuZWFzeWljb24uY24vZmlsZXMvNTcvNTc5MjYzLnBuZ6Sjum4AAAAASUVORK5CYII=) ![](https://img.shields.io/badge/License-AGPL--3.0_License-yellow?style=flat) ![](https://img.shields.io/badge/Version-Basic-blueviolet?style=flat)

## Start

Clone/Download Repositories | 下载/拉取代码

`$ git clone https://github.com/0ojixueseno0/go-Pichubot.git`



## Example

```go
// main.go
//...
// 消息事件
case "message":
  sender := receive["sender"].(map[string]interface{})
  switch receive["message_type"] {
  // 私聊信息
  case "private":
   -> Code Here <-
//...

```
