package modules

import (
	"fmt"
	"github.com/alsm/forecastio"
	"github.com/brambu/brambu-telegram-bot/config"
	. "github.com/brambu/brambu-telegram-bot/helpers"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type Weather struct {
	config   config.BotConfiguration
	client   *forecastio.APIConn
	geocoder geo.Geocoder
}

func (w *Weather) Name() *string {
	name := "weather"
	return &name
}

func (w *Weather) LoadConfig(conf config.BotConfiguration) {
	w.config = conf
	client := forecastio.NewConnection(w.config.DarkskyToken)
	err := client.SetUnits("auto")
	if err != nil {
		log.Error().Err(err).Msg("weather darksky set units error")
	}
	w.client = client
	w.geocoder = openstreetmap.Geocoder()
}

func (w *Weather) Evaluate(update tgbotapi.Update) bool {
	return CheckPrefix(update, "/weather")
}

func (w *Weather) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Info().Msg("Sending weather.")
	searchText := strings.Join(strings.Split(GetUpdateMessageText(update), " ")[1:], " ")
	location := w.getLocation(searchText)
	if location == nil {
		ReplyWithText(bot, update, "aroo?")
		return
	}
	ReplyWithText(bot, update, w.getWeather(location))
}

func (w *Weather) getWeather(location *geo.Location) string {
	address := w.getAddress(location)
	f, err := w.client.Forecast(location.Lat, location.Lng, []string{}, false)
	log.Info().
		Int("api_calls", w.client.APICalls()).
		Msg("weather darksky api call counter")
	if err != nil {
		return "aroo?"
	}
	f.ParseTimes()
	u := "C"
	wu := "mps"
	switch f.Flags.Units {
	case "us":
		u = "F"
		wu = "mph"
	case "ca":
		wu = "kph"
	case "uk2":
		wu = "mph"
	}
	t := timeIn(f.Currently.Time, f.Timezone)
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

func (w *Weather) getLocation(searchString string) *geo.Location {
	res, err := w.geocoder.Geocode(searchString)
	if err != nil {
		log.Error().Err(err).
			Str("search_string", searchString).
			Msg("weather error getting location")
	}
	return res
}

func (w *Weather) getAddress(location *geo.Location) *geo.Address {
	res, err := w.geocoder.ReverseGeocode(location.Lat, location.Lng)
	if err != nil {
		log.Error().Err(err).
			Msg("weather error getting address")
	}
	return res
}

func timeIn(t time.Time, name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t
}
