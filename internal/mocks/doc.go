/*
Package mocks will have all the mocks of the app.
*/
package mocks

//go:generate mockery -case underscore -output ./forward -dir ../forward -name Notifier
//go:generate mockery -case underscore -output ./forward -dir ../forward -name Service
//go:generate mockery -case underscore -output ./deadmansswitch -dir ../deadmansswitch -name Service

//go:generate mockery -case underscore -output ./notify/telegram -dir ../notify/telegram -name Client
//go:generate mockery -case underscore -output ./notify -dir ../notify -name TemplateRenderer
