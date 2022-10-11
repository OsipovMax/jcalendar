# jcalendar

### Требования

Сервис должен иметь HTTP API, позволяющее:
- [x] Cоздать пользователя;
- [x] Cоздать встречу в календаре пользователя со списком приглашенных пользователей;
- [x] Получить детали встречи;
- [x] Принять или отклонить приглашение другого пользователя;
- [ ] Найти все встречи пользователя для заданного промежутка времени;
- [ ] Для заданного списка пользователей и минимальной продолжительности встречи, найти ближайшей интервал времени, в котором все эти пользователи свободны;
- [ ] У встреч в календаре должна быть возможна настройка повторов (В повторах нужно поддержать все возможности, доступные в Google-календаре, кроме Сustom.).

**Необязательные требования:**
- [x] Аутентификация пользователя;
- [x] Поддержка видимости встреч (если встреча приватная, другие пользователи могут получить только информацию о занятости пользователя, но не детали встречи);
- [ ] Настройки часового пояса пользователя и его рабочего времени, использование этих настроек для поиска интервала времени, в котором участники свободны;
- [ ] Настройки нотификации пользователя перед встречей (саму нотификацию достаточно реализовать записью в лог);
- [ ] Поддержка Custom повторов, как в Google-календаре;