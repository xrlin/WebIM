import {generateURL} from "../utils/url";
import request from "../utils/request";

export function retrieveUnreaOfflinedMessages() {
  let options = {
    method: 'GET'
  };
  let url = generateURL('/api/messages/unread');
  return request(url, options)
}

export function ackMessages(messageIds) {
  let options = {
    method: 'DELETE',
    body: JSON.stringify({message_ids: messageIds})
  };
  let url = generateURL('/api/messages/ack');
  return request(url, options)
}
