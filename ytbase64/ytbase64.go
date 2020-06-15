package ytbase64

import "log"

func Encode64(int_id int64) string {
	var str_id string
	var temp uint8 = uint8(int_id&0xF) << 2
	int_id >>= 4
	for i := 0; i < 11; i++ {
		switch {
		case temp < 26:
			temp += 65
		case temp < 52:
			temp += 71
		case temp < 62:
			temp -= 4
		case temp == 62:
			temp = '-'
		case temp == 63:
			temp = '_'
		default:
			log.Fatal("Error Encode64: unknown character\n")
		}

		str_id = string(temp) + str_id
		temp = uint8(int_id & 63)
		int_id >>= 6
	}
	return str_id
}

func Decode64(str_id string) int64 {
	var int_id int64
	var temp int64 = 0

	for i := 0; i < 11; i++ {
		temp = int64(str_id[10-i])
		switch {
		case temp >= 'A' && temp <= 'Z':
			temp -= 65
		case temp >= 'a' && temp <= 'z':
			temp -= 71
		case temp >= '0' && temp <= '9':
			temp += 4
		case temp == '-':
			temp = 62
		case temp == '_':
			temp = 63
		default:
			log.Fatal("Error Decode64: unknown character\n")
		}
		if i != 0 {
			temp = temp << (i*6 - 2)
		} else {
			temp = temp >> 2
		}
		int_id |= temp
	}
	return int_id
}
