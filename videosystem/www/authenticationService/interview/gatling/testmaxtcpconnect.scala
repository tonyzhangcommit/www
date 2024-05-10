import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class TestMaxTcpConnect extends Simulation {

  val httpProtocol = http
    .baseUrl("http://39.105.9.252:9999/api") // 请求根路径以及设置请求头



  val scn = scenario("TestMaxTcpConnect") // A scenario is a chain of requests and pauses
    .exec(http("testtcpconnect")
    .get("/test"))

  // 组装请求
  setUp(
    scn.inject(
      atOnceUsers(4000)
      // constantUsersPerSec(800) during (60 seconds) // Inject 1000 users per second for a duration of 10 seconds
    ).protocols(httpProtocol)
  )
}
