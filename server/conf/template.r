args<-commandArgs(T)
{x_data}
y<-c({y_data})
fit<-lm(y~{x_expr})
summary(fit)
png(file = "{workspace}/rp%03d.png")
plot(fit)
q("no")