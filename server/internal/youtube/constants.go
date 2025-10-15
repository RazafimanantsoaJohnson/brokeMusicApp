package youtube

const ytDlpGreedyJsonRegex = `\{.+\}`
const youtubeBaseUrl = "https://www.youtube.com/watch"
const youtubeAPIBaseURl = "https://www.googleapis.com/youtube/v3/search"

// const ytDlpJsonRegex = `\{.+?\}`
const ytInitialPlayerResponseRegexString = `var ytInitialPlayerResponse\s*=\s*(\{.+?\});` // `var ytInitialPlayerResponse\s*=`
const ytInitialDataRegexString = `(?:window\s*\[\s*["\']ytInitialData["\']\s*\]|ytInitialData)\s*=\s*(\{.+?\});`
const ytCfgRegexString = `ytcfg\.set\s*\(\s*({.+?})\s*\)\s*;`
const ytPlayerBaseUrl = "https://youtube.com"
const downloadSuccessRegex = `.*\[download\]\s+100.*`
