# Небольшие заметки для себя
## Viper - библиотека для конфигурации
## Cobra - библиотека для создания CLI

## [Viper] - Struct Tags
```
type Example struct {
    // Basic mapstructure tag
    Field1 string `mapstructure:"field_1"`

    // Squash embeds the fields at the same level
    Field2 Config `mapstructure:",squash"`

    // Omits the field if empty
    Field3 string `mapstructure:"field_3,omitempty"`

    // Remaining tags for special cases
    Field4 string `mapstructure:"field_4,remain"`
    Field5 string `mapstructure:"field_5,omitzero"`
    Field6 string `mapstructure:"field_6,dive"`
}
```
Description:
* `squash` - вложенные структуры будут сжаты в один уровень
* `omitempty` - не будет создаваться поле, если оно пустое
* `remain` - остальные поля будут сохранены
* `omitzero` - не будет создаваться поле, если оно равно нулю
* `dive` - вложенные структуры будут расширены
* `squash,remain` - вложенные структуры будут сжаты в один уровень, а остальные поля будут сохранены
* `squash,omitempty` - вложенные структуры будут сжаты в один уровень, а пустые поля будут пропущены
* `mapstructure:"name"` - basic field mapping
* `mapstructure:",squash"` - flattens nested structs
* `mapstructure:"name,omitempty"` - skips empty values
* `mapstructure:"name,remain"` - keeps unmapped values
* `mapstructure:"name,omitzero"` - omits zero values
* `mapstructure:"name,dive"` - for slice/map element processing

## Magic date format GO
1. Магическая дата: Go использует конкретную дату и время как образец - 15:04:05 2 января 2006 года, часовой пояс -0700 MST. Это соответствует числам 1, 2, 3, 4, 5, 6, 7 (месяц, день, час, минута, секунда, год, часовой пояс).


2. Мнемоническое правило: 01/02 03:04:05PM '06 -0700

    01 = месяц (January) \
    02 = день \
    03 = час (12-часовой формат) \
    04 = минута \
    05 = секунда \
    06 = год (2006) \
    -0700 = часовой пояс 


3. Шпаргалка для распространенных форматов:

    "2006-01-02" для ISO даты (YYYY-MM-DD) \
    "01/02/2006" для американского формата (MM/DD/YYYY) \
    "02/01/2006" для европейского формата (DD/MM/YYYY) \
    "15:04:05" для времени в 24-часовом формате \
    "3:04 PM" для времени в 12-часовом формате 