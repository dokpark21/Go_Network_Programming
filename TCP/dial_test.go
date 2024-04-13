package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buffer := make([]byte, 1024)
				for {
					n, err := c.Read(buffer)
					if err != nil {
						if err != io.EOF {
							t.Log(err)
						}
						return
					}
					t.Logf("received: %q", buffer[:n])
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}

	conn.Close()
	<-done
	listener.Close()
	<-done
}

// 작동 순서
/*
net.Listen 호출: TestDial 테스트 함수가 시작되면서, net.Listen 함수를 호출하여 TCP 네트워크의 127.0.0.1:0 주소에 대한 리스너를 생성합니다. 이때, 0 포트는 시스템이 사용 가능한 포트를 자동으로 선택하게 합니다.

첫 번째 고루틴 시작: go func() {...}를 사용하여 첫 번째 비동기 고루틴을 시작합니다. 이 고루틴은 연결을 수락(Accept)하고, 각 연결에 대해 읽기 작업을 수행하는 또 다른 고루틴을 시작합니다.

Accept 루프: 첫 번째 고루틴은 무한 루프에서 listener.Accept()를 호출하여 들어오는 연결을 기다립니다. 새로운 연결이 들어올 때마다 다음 단계로 진행합니다.

새로운 연결 처리를 위한 두 번째 고루틴 시작: 새 연결이 들어오면, Accept가 새 연결(conn)을 반환하고, 이 연결을 처리하기 위해 새로운 고루틴을 시작합니다(go func(c net.Conn) {...}(conn)).

연결 읽기: 두 번째 고루틴은 연결에서 데이터를 읽습니다(c.Read(buffer)). 데이터가 수신되면, "received: %q" 로그를 출력합니다. EOF(End Of File)에 도달하거나 오류가 발생하면 연결을 닫고 고루틴을 종료합니다.

net.Dial 호출: 메인 테스트 고루틴에서 net.Dial을 호출하여 리스너 주소로 연결을 시도합니다. 이는 위에서 생성한 리스너에 대한 연결을 만듭니다.

연결 닫기: conn.Close()를 호출하여 연결을 즉시 닫습니다.

고루틴 종료 신호 대기: 메인 테스트 고루틴은 <-done을 통해 고루틴의 종료 신호를 기다립니다. 첫 번째 고루틴과 두 번째 고루틴 모두 done 채널에 신호를 보내야 합니다.

리스너 닫기: listener.Close()를 호출하여 더 이상의 들어오는 연결을 수락하지 않습니다.

테스트 종료: 모든 고루틴의 실행이 완료되고, 모든 리소스가 정리되면 테스트는 종료됩니다.
*/

// 테스트 코드 로직
/*
Listener 설정: 코드는 net.Listen을 사용하여 TCP 네트워크의 특정 주소(127.0.0.1:0)에서 리스너를 생성합니다. 여기서 :0은 시스템이 자동으로 사용 가능한 포트를 선택하게 합니다.

고루틴과 연결 처리: 리스너가 설정되면, 첫 번째 고루틴이 시작되어 들어오는 연결을 대기하고 Accept를 통해 수락합니다. 새 연결이 성립되면, 이 연결을 처리하기 위해 새로운(두 번째) 고루틴이 생성됩니다.

net.Dial을 통한 연결 시도: 메인 테스트 고루틴에서는 net.Dial을 사용해 리스너에 연결을 시도합니다. 이 연결이 성공적으로 이루어지면, 리스너의 Accept 호출이 반환되고, 새로운 연결을 처리하기 위한 두 번째 고루틴이 시작됩니다.

데이터 전송 없음과 연결 종료: 맞습니다, 이 테스트 코드에서는 net.Dial을 통해 연결한 후 아무런 데이터를 전송하지 않고 바로 연결을 종료합니다(conn.Close()). 연결이 종료되면, 연결을 처리하고 있던 두 번째 고루틴이 EOF를 감지하고 종료합니다.

고루틴 종료 신호와 메인 고루틴의 대기: 각 고루틴은 종료 시점에 done 채널에 신호(빈 구조체 {})를 보냅니다. 메인 테스트 고루틴은 이 신호를 두 번 받을 때까지 대기합니다(<-done). 첫 번째 신호는 연결을 처리하는 고루틴의 종료를 나타내고, 두 번째 신호는 Accept 루프를 처리하는 첫 번째 고루틴의 종료를 나타냅니다.

리소스 정리와 테스트 종료: 모든 연결이 종료되고, 모든 고루틴의 작업이 완료된 후, 리스너는 닫히고(listener.Close()), 테스트는 종료됩니다.
*/
