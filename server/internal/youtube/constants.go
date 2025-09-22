package youtube

const ytInitialPlayerResponseRegexString = `var ytInitialPlayerResponse\s*=\s*(\{.+?\});` // `var ytInitialPlayerResponse\s*=`
const ytInitialDataRegexString = `(?:window\s*\[\s*["\']ytInitialData["\']\s*\]|ytInitialData)\s*=\s*(\{.+?\});`
const ytCfgRegexString = `ytcfg\.set\s*\(\s*({.+?})\s*\)\s*;`
const ytPlayerBaseUrl = "https://youtube.com"
