import {generateURL} from "../utils/url";
import request from "../utils/request";

export function leaveRoom(roomID) {
  let options = {
    method: 'DELETE',
  };
  let url = generateURL(`/api/rooms/${ roomID }/leave`);
  return request(url, options)
}
