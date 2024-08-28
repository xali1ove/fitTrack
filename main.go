package main

import (
	"fmt"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре
)

// Training общая структура для всех тренировок
// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // тип тренировки
	Action       int           // количество шагов/гребков
	LenStep      float64       // длина одного шага/гребка в метрах
	Duration     time.Duration // продолжительность тренировки
	Weight       float64       // вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	return t.distance() / t.Duration.Seconds() * 3600
}

// Calories возвращает количество потраченных килокалорий на тренировке.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потраченных килокалорий при беге.
const (
	CaloriesMeanSpeedMultiplier = 18   // множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // коэффициент изменения средней скорости
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возвращает количество потраченных килокалорий при беге.
func (r Running) Calories() float64 {
	speed := r.meanSpeed()
	return (CaloriesMeanSpeedMultiplier*speed + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Seconds() / 3600
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (r Running) TrainingInfo() InfoMessage {
	info := r.Training.TrainingInfo()
	info.Calories = r.Calories()
	return info
}

// Константы для расчета потраченных килокалорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // коэффициент для роста
	KmHInMsec                     = 0.278 // коэффициент для перевода км/ч в м/с
)

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64 // рост пользователя в сантиметрах
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
func (w Walking) Calories() float64 {
	speed := w.meanSpeed() * KmHInMsec
	heightInMeters := w.Height / CmInM
	return (CaloriesWeightMultiplier*w.Weight + (speed*speed/heightInMeters)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Seconds() / 3600
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (w Walking) TrainingInfo() InfoMessage {
	info := w.Training.TrainingInfo()
	info.Calories = w.Calories()
	return info
}

// Константы для расчета потраченных килокалорий при плавании.
const (
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // множитель веса пользователя
)

// Swimming структура, описывающая тренировку Плавание
// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool int // длина бассейна в метрах
	CountPool  int // количество пересечений бассейна
}

// meanSpeed возвращает среднюю скорость при плавании.
func (s Swimming) meanSpeed() float64 {
	return float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Seconds() * 3600
}

// Calories возвращает количество калорий, потраченных при плавании.
func (s Swimming) Calories() float64 {
	speed := s.meanSpeed()
	return (speed + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Seconds() / 3600
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     s.distance(),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	return fmt.Sprint(info)
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
