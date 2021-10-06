(async function () {
    const year = 2021
    const month = 4

    const pageNum = 1
    const pageSize = 100
    const u = `https://api.bilibili.com/pgc/season/index/result?season_version=-1&spoken_language_type=1&area=2&is_finish=-1&copyright=-1&season_status=-1&season_month=${month}&year=%5B${year}%2C${year + 1})&style_id=-1&order=3&st=1&sort=0&page=${pageNum}&season_type=1&pagesize=${pageSize}&type=1`
    console.log(u)
    const resp = await fetch(u)
    const j = await resp.json()

    if (j.data.total > pageSize) throw '数量过多：' + j.data.total
    const list = j.data.list
    const no = (t, s) => t.indexOf(s) === -1
    const res = list.map(({cover, title, link}) => {
        return {
            cover: cover.replace('http://', 'https://'),
            title: title,
            link: link.replace('http://', 'https://'),

            yearMonth: year * 100 + month,
            episode: '',

            markStatus: 0,  // 0:还没看 -1:放弃 1:追番中 2:已追完
            markScore: 0,
            markBrev: '',
            markDate: '',
            markEpisode: '',
            tags: [],
        }
    }).filter(({title}) => no(title, '中配版') && no(title, '粤配版') && no(title, '地區'))
    console.log(JSON.stringify(res))
})()
