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