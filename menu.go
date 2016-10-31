package fbbot

import ()

type Menu struct {
	Type                string // Value is web_url or postback, required
	Title               string // Button title, required, has a 30 character limit
	URL                 string // For web_url buttons, this URL is opened in a mobile browser when the button is tapped. Required if type is web_url
	Payload             string // For postback buttons, this data will be sent back to you via webhook. Required if type is postback. Has a 1000 character limit
	WebviewHeightRatio  string // Valid values: compact, tall, full
	MessengerExtensions bool   // Must be true if using Messenger Extensions. https://developers.facebook.com/docs/messenger-platform/send-api-reference/webview
}

// Persistent Menu

// curl -X POST -H "Content-Type: application/json" -d '{
//   "setting_type" : "call_to_actions",
//   "thread_state" : "existing_thread",
//   "call_to_actions":[
//     {
//       "type":"postback",
//       "title":"Help",
//       "payload":"DEVELOPER_DEFINED_PAYLOAD_FOR_HELP"
//     },
//     {
//       "type":"postback",
//       "title":"Start a New Order",
//       "payload":"DEVELOPER_DEFINED_PAYLOAD_FOR_START_ORDER"
//     },
//     {
//       "type":"web_url",
//       "title":"Checkout",
//       "url":"http://petersapparel.parseapp.com/checkout",
//       "webview_height_ratio": "full",
//       "messenger_extensions": true
//     },
//     {
//       "type":"web_url",
//       "title":"View Website",
//       "url":"http://petersapparel.parseapp.com/"
//     }
//   ]
// }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=PAGE_ACCESS_TOKEN"

// call_to_actions is limited to 5

// Success
// {
//   "result": "Successfully added new_thread's CTAs"
// }

// Delete

// curl -X DELETE -H "Content-Type: application/json" -d '{
//   "setting_type":"call_to_actions",
//   "thread_state":"existing_thread"
// }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=PAGE_ACCESS_TOKEN"
