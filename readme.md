# Instagram Notification Service

## Setup

First of all you need to setup google service account for your Google Sheet

[Documentation of setting up a Google Sheet](https://docs.google.com/document/d/1wEXHRuFtRBwR41zBkk8cOKW9BkhJeHqAy32zyX8p5hk/edit?usp=sharing)

After that rename your .env.example to .env and fill your env values
``cp .env.example .env``


So that's it! Then run docker build in root of project and all things will be done!

## Documentation

[Trello board card](https://trello.com/c/1l4oROzH/316-%D1%80%D0%B0%D1%81%D1%81%D1%8B%D0%BB%D0%BA%D0%B0-%D0%B2-%D0%B4%D0%B8%D1%80%D0%B5%D0%BA%D1%82)

All documentation about endpoins you can find by route
``GET /swagger/``

Or in ``./docs`` folder

If you change routes you MUST run ``make swag`` for generating swagger documentation

## Testing

For testing run mocks
``make gen``

After that you can run test
``make test``