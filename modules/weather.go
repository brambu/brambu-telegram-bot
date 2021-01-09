package modules

import (
	"fmt"
	"github.com/alsm/forecastio"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

type Weather struct {
	config config.BotConfiguration
}

func (w *Weather) LoadConfig(conf config.BotConfiguration) {
	w.config = conf
}

func (w Weather) GetLocation(searchString string) *geo.Location {
	g := openstreetmap.Geocoder()
	res, err := g.Geocode(searchString)
	if err != nil {
		log.Printf("Weather error getting location: %s, %s", searchString, err)
	}
	return res
}

func (w Weather) GetAddress(location *geo.Location) *geo.Address {
	g := openstreetmap.Geocoder()
	res, err := g.ReverseGeocode(location.Lat, location.Lng)
	if err != nil {
		log.Printf("Weather error getting address: %s", err)
	}
	return res
}

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func (w Weather) GetWeather(location *geo.Location) string {
	address := w.GetAddress(location)
	c := forecastio.NewConnection(w.config.DarkskyToken)
	err := c.SetUnits("auto")
	if err != nil {
		log.Printf("Weather Darksky set units error: %s", err)
	}
	f, err := c.Forecast(location.Lat, location.Lng, []string {}, false)
	if err != nil {
		return "aroo?"
	}
	f.ParseTimes()
	u := "C"
	wu := "mps"
	switch {
	case f.Flags.Units == "us":
		u = "F"
		wu = "mph"
	case f.Flags.Units == "ca":
		wu = "kph"
	case f.Flags.Units == "uk2":
		wu = "mph"
	}
	t, _ := TimeIn(f.Currently.Time, f.Timezone)
	log.Printf("Weather Darksky API Calls made today: %d\n", c.APICalls())
	retSlice := []string{
		fmt.Sprintf("Current Weather for %s %s %s at %s\n",
			address.City, address.State, address.CountryCode, t.Format("Jan 02, 2006 15:04")),
		fmt.Sprintf("_%s_ _%s_ _%s_\n", f.Minutely.Summary, f.Hourly.Summary, f.Daily.Summary),
		fmt.Sprintf("Temperature: *%.0f°%s*", f.Currently.Temperature, u),
		fmt.Sprintf("Wind: %.0f%s  Humidity %.0f%%", f.Currently.WindSpeed, wu, f.Currently.Humidity * 100),
		fmt.Sprintf("High: %.0f°%s Low: %.0f°%s",
			f.Daily.Data[0].TemperatureMax, u, f.Daily.Data[0].TemperatureMin, u),
	}
	for _, alert := range f.Alerts {
		alertsSlice := []string{
			fmt.Sprintf("\n*Alert*: [%s](%s)", alert.Title, alert.URI),
		}
		retSlice = append(retSlice, alertsSlice...)
	}
	return strings.Join(retSlice, "\n")
}

func (w Weather) Evaluate(update tgbotapi.Update) bool {
	if strings.HasPrefix(strings.ToLower(update.Message.Text), "/weather") {
		log.Printf("Weather command from [%d,%s]: %s",
			update.Message.From.ID,
			update.Message.From.UserName,
			update.Message.Text)
		return true
	}
	return false
}

func (w Weather) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Println("Sending weather.")
	searchText := strings.Join(strings.Split(update.Message.Text, " ")[1:], " ")
	location := w.GetLocation(searchText)
	if location == nil {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "aroo?")
		_, err := bot.Send(message)
		if err != nil {
			log.Printf("Warning: Weather nolocation error %s", err)
		}
		return
	}
	weather := w.GetWeather(location)

	message := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
	message.ParseMode = "Markdown"
	message.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(message)
	if err != nil {
		log.Printf("Warning: Weather error %s", err)
	}
}
