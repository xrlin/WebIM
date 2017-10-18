import {getImageUrls} from "./url";

export function htmlToMessage(html) {
  let result = html.replace(/<br>/g, '\n');
  result = result.replace(/<img.*? data-faceid="(.*?)".*? data-facetext="(.*?)".*?>/gi, '[$2_$1]');
  return result
}

export function messagetToHtml(message) {
  let content = message.content;
  let result = '';
  if (message.msg_type === 6) {
    result = `<iframe frameborder="no" border="0" marginwidth="0" marginheight="0" width=330 height=86 src="//music.163.com/outchain/player?type=2&id=${content}&auto=0&height=66"></iframe>`
  } else if (message.msg_type === 4) {
    let {thumbnailUrl} = getImageUrls(content);
    result = `<img src="${thumbnailUrl}" alt="图片不存在" >`
  } else {
    result = content.replace(/\[(.*?)_(.*?)]/gi, '<img src="/$2.gif">');

  }
  return result
}
