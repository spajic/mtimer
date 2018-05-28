# mtimer

`mtimer` позволяет сохранить `mtime` для файлов в директории и потом восстановить сохранённые значения.

`git` принципиально не сохраняет время изменения контроллируемых файлов. В случае, если время изменения файла важно (например, для кэширования), `mtimer` поможет их восстановить.

`mtimer` написан на `go` и работает быстро.

## usage

### Сохранить mtimes в файл mtimer.dat

`mtimer --store --filespath=/path/to/files --timespath=/path/to_mtimer_dat --ignore=node_modules,tmp,.git`

- `--store` - режим сохранения mtimes в файл
- `--filespath` - путь к директории, для файлов которой нужно **сохранить** mtimes
- `--timespath` - путь к директории, куда сохранить файл `mtimes.dat`
- `--ignore` - список поддиректорий, для которых mtimes сохранять не нужно

### Восстановить mtimes из файла mtimer.dat

`mtimer --apply --filespath=/path/to/files --timespath=/path/to_mtimer_dat`

- `--apply` - режим восстановления mtimes из файла
- `--filespath` - путь к директории, для файлов которой нужно **восстановить** mtimes
- `--timespath` - путь к директории, откуда взять файл `mtimes.dat`

### Показать справку

`mtimer --help`

### Показать версию

`mtimer --version`

## comparison

Попробовал подобное решение на `perl` - https://github.com/danny0838/git-store-meta - в нашем проекте работало порядка 10 секунд, `mtimer` же отрабатывает мгновенно.