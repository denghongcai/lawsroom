#!/bin/bash

sed -ri 's#https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.6/styles/default.min.css#https://dn-txthinking.qbox.me/static/highlight-8.6.min.css#' themes/material-design/layouts/partials/header.html
sed -ri 's#https://code.jquery.com/jquery-2.1.4.min.js#https://dn-txthinking.qbox.me/static/jquery-2.1.4.min.js#' themes/material-design/layouts/partials/footer.html
sed -ri 's#https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.6/highlight.min.js#https://dn-txthinking.qbox.me/static/highlight-8.6.min.js#' themes/material-design/layouts/partials/footer.html

