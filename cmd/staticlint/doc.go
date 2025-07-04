// Package staticlint запускает набор статических анализаторов:
// - стандартные анализаторы пакета go/analysis;
// - все SA-анализаторы staticcheck (например, SA1000 — неверное использование time.Parse);
// - ST1000 — проверка имени пакета;
// - собственный noosexit — запрещает использовать os.Exit в main.
//
// Запуск:
//
//	go run ./cmd/staticlint ./...
package main
