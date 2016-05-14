#!/bin/bash

sed -ri 's#href="http.*?Roboto\+Mono.*?"#"/dist/css/RobotoMono.css"#' bower_components/font-roboto/roboto.html
sed -ri 's#href="http.*?Roboto.*?"#"/dist/css/Roboto.css"#' bower_components/font-roboto/roboto.html
