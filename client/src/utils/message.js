export function htmlToMessage(html) {
  let result = html.replace(/<br>/g, '\n');
  result = result.replace(/<img.*? data-faceid="(.*?)".*? data-facetext="(.*?)".*?>/gi, '[$2_$1]');
  return result
}

export function messagetToHtml(messageText) {
  let result = messageText.replace(/\[(.*?)_(.*?)]/gi, '<img src="/$2.gif">');
  return result
}
