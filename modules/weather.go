package modules

import (
	"fmt"
	"github.com/alsm/forecastio"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).
			Str("search_string", searchString).
			Msg("weather error getting location")
	}
	return res
}

func (w Weather) GetAddress(location *geo.Location) *geo.Address {
	g := openstreetmap.Geocoder()
	res, err := g.ReverseGeocode(location.Lat, location.Lng)
	if err != nil {
		log.Error().Err(err).
			Msg("weather error getting address")
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
		log.Error().Err(err).Msg("weather darksky set units error")
	}
	f, err := c.Forecast(location.Lat, location.Lng, []string{}, false)
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
	log.Info().
		Int("api_calls", c.APICalls()).
		Msg("weather darksky api calls made today")
	retSlice := []string{
		fmt.Sprintf("Current Weather for %s %s %s at %s\n",
			address.City, address.State, address.CountryCode, t.Format("Jan 02, 2006 15:04")),
		fmt.Sprintf("_%s_ _%s_ _%s_\n", f.Minutely.Summary, f.Hourly.Summary, f.Daily.Summary),
		fmt.Sprintf("Temperature: *%.0f°%s*", f.Currently.Temperature, u),
		fmt.Sprintf("Wind: %.0f%s  Humidity %.0f%%", f.Currently.WindSpeed, wu, f.Currently.Humidity*100),
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
		log.Info().
			Int("from_id", update.Message.From.ID).
			Str("from_user_name", update.Message.From.UserName).
			Str("text", update.Message.Text).
			Msg("weather command")
		return true
	}
	return false
}

func (w Weather) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Info().Msg("Sending weather.")
	searchText := strings.Join(strings.Split(update.Message.Text, " ")[1:], " ")
	location := w.GetLocation(searchText)
	if location == nil {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "aroo?")
		_, err := bot.Send(message)
		if err != nil {
			log.Error().Err(err).Msg("weather nolocation error")
		}
		return
	}
	weather := w.GetWeather(location)

	message := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
	message.ParseMode = "Markdown"
	message.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("weather error sending message")
	}
}
