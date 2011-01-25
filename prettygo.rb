#!/usr/bin/env ruby1.9.1

require 'rubygems'
require 'termcolor'
require 'open4'


RULES = {
	/make.*\n/ => "",
	/rm.*\n/ => "",
	/cp.*\n/ => "",
	/gotest\n/ => "",
	/gopack.*\n/ => "",
	/(8|5|6)g -o _.+_\.(8|5|6) (?<files>.+)\n/ => "<green>Compiling:</green> \\k<files>\n",
	/(?<error>\w+\.go:\d+: .+)\n/ => "<red>\\k<error></red>\n"
}

Open4::open4("gomake") do |cid, i, o, e|
	while true do
		begin
			line = o.readline
			RULES.each_key do |k|
				if line.match(k) then 
					line.gsub!(k, RULES[k])
					break
				end
			end
			print line.termcolor
		rescue
			puts "Done!"
			break
		end
	end
end


