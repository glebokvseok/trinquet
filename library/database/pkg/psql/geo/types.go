package geo

import (
	"database/sql/driver"
	"fmt"
	"github.com/cridenour/go-postgis"
	"github.com/pkg/errors"
)

type Point struct {
	Longitude float64
	Latitude  float64
}

func (point *Point) Value() (driver.Value, error) {
	if point == nil {
		return nil, nil
	}

	return fmt.Sprintf("SRID=4326;POINT(%f %f)", point.Longitude, point.Latitude), nil
}

func (point *Point) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.Errorf("unsupported type %T for Point.Scan", value)
	}

	var pgPoint postgis.PointS
	if err := pgPoint.Scan(data); err != nil {
		return errors.WithStack(err)
	}

	point.Longitude, point.Latitude = pgPoint.X, pgPoint.Y

	return nil
}
