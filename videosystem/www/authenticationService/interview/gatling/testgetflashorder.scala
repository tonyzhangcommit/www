import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class TestGetFlashEventOrder extends Simulation {

  val httpProtocol = http
    .baseUrl("http://39.105.9.252:9999/api/flashevent") // 请求根路径以及设置请求头
    .contentTypeHeader("application/json")

  // 从json中读取数据  circular  表示到文件尾后会重头开始读取
  val feeder = jsonFile("user-files/data/token.json").circular

    
  val scn = scenario("TestGetFlashEventOrder") // 设置了请求的一个动作，比如请求路径，请求方式等
  .feed(feeder)
    .exec(
      http("takeorder")   // 本次请求的名称
      .post("/takeorder")  // 具体请求路径
      .header("Authorization", "Bearer ${utoken}")
      .body(StringBody(
      """
      {"eventid":22,"peoductid":1,"userid":${userid},"count":1}
      """
    )).asJson
    .check(status.is(200))
    ) 

  // 组装请求
  setUp(
    scn.inject(
      atOnceUsers(4000),                   // 初始 4000 用户立即请求
      nothingFor(3 seconds),               // 继续等待 3 秒
      atOnceUsers(400),                    // 再注入 400 用户
      nothingFor(3 seconds),               // 再等待 3 秒
      atOnceUsers(40)  
      // constantUsersPerSec(800) during (60 seconds) // Inject 1000 users per second for a duration of 10 seconds
    ).protocols(httpProtocol)
  )
}
