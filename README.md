###2022.02.11

## Błędy znalezione po klikaniu frontendu
1. W repozytorium nie trzymamy plików .env i .db oraz dziwnych .exe jak ten swag.exe
2. Brakuje mi errorów na stronie frontendowej, żebym widział co się dzieje bez wchodzenia w Network czy Consolę
2.1 Jeżeli jakaś rola nie ma dostępów do danej opcji (np. Dodania studenta), to powinno być to ukryte na frontendzie, dopóki ktoś się nie zaloguje z odpowiednią rolą.
3. Po dodaniu studenta muszę odświeżyć stronę żeby go zobaczyć
4. Jeżeli wpiszę ID studenta, które nie istnieje, to i tak dostanę 200 z backendu.
4.1 Jeżeli zedytuję usera, którego nie ma w bazie, to mi go automatycznie stworzy, tak nie powinno być.
5. Jeżeli przy updacie użytkownika podam tylko jedną kolumnę do aktualizacji np. NAME, to z backendu dostanę 200, a nic się nie zmieni w bazie
5.1 POPRAWKA: Nie działa updatowanie usera.

## Swagger
5. Nie ma możliwości zautentykowania się do endpointów na swaggerze.
5.1 Swagger jest bardzo słabo opisany, nie ma przykładów responsów, jak i errorów 
6. Przy dodawaniu usera na Swaggerze wymagane są wszystkie pola, a na stronie frontendu mogę sobie podać tylko jedno i i tak request przejdzie