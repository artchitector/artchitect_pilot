export default function ({$axios, store}) {
  // if (process.server) {
  //   return
  // }
  $axios.onRequest((req) => {
    console.log('axios onRequest', req)
    req.headers.common['Authorization'] = "Bearer: abcdef"
  })
}
