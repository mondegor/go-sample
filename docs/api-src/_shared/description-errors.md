  ### Пользовательские ограничения полей и ошибки

  #### Ограничения полей
  - required - поле обязательно для заполнения;
  - omitempty - поле может быть не указано (не будет использоваться методом, в который было передано);
  - min=N - поле должно быть не менее N символов;
  - max=N - поле должно быть не более N символов;
  - gte=N - числовое поле должно быть равно или более N;
  - lte=N - числовое поле должно быть равно или менее N;
  - enum - поле должно содержать одно из ENUM значений;
  - pattern=P - поле должно соответствовать регулярному выражению P;

  #### Ошибки
  - ErrVersionInvalid - если передаваемая версия объекта не совпала с текущей версией объекта.\
    Как правило, это означает, что объект был ранее изменён другим процессом;
  - ErrSwitchStatusRejected - перевод в указанный статус объекта отклонён.\
    WorkFlow объекта запрещает переключение в указанный статус;