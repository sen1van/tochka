#include <WiFi.h>

const char* ssid = "p13-2Ghz";
const char* password = "p13123456";

// Указываем сервер и порт правильными типами данных
const char* server_host = "sen1van.ru";
const uint16_t server_port = 8080;

WiFiClient client;

void setup() {
  Serial.begin(115200);
  
  WiFi.begin(ssid, password);
  Serial.print("Подключение к Wi-Fi...");

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print("."); // Просто выводим точки, чтобы видеть прогресс
  }

  Serial.println("\nГотово! Wi-Fi подключен.");
  Serial.print("IP-адрес платы: ");
  Serial.println(WiFi.localIP());
}

void loop() {
  // Проверяем, ввели ли мы что-то в Монитор порта
  if (Serial.available() > 0) {
    String text = Serial.readString();
    text.trim(); // Убираем лишние пробелы и знаки переноса строки

    Serial.print("Отправка запроса на путь: ");
    Serial.println(text);

    // Подключаемся к серверу КАЖДЫЙ РАЗ перед отправкой запроса
    if (client.connect(server_host, server_port)) {
      Serial.println("Успешно подключились к серверу!");

      // Отправляем HTTP GET-запрос
      client.println("GET " + text + " HTTP/1.1");
      client.println("Host: " + String(server_host) + ":" + String(server_port));
      client.println("Connection: close");
      client.println(); // Пустая строка — конец запроса

      // Ожидаем ответа от сервера (таймаут 5 секунд)
      unsigned long timeout = millis();
      while (client.available() == 0) {
        if (millis() - timeout > 5000) {
          Serial.println(">>> Ошибка: Сервер не ответил вовремя!");
          client.stop();
          return;
        }
      }

      // Читаем и выводим ответ от сервера в Монитор порта
      while (client.available()) {
        char c = client.read();
        Serial.write(c);
      }
      
      // Закрываем соединение, так как мы указали Connection: close
      client.stop();
      Serial.println("\n--- Соединение закрыто ---");

    } else {
      Serial.println("Ошибка подключения к серверу");
    }
  }
}
