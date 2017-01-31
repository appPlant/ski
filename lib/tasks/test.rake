#
# Copyright (c) 2013-2016 by appPlant GmbH. All rights reserved.
#
# @APPPLANT_LICENSE_HEADER_START@
#
# This file contains Original Code and/or Modifications of Original Code
# as defined in and that are subject to the Apache License
# Version 2.0 (the 'License'). You may not use this file except in
# compliance with the License. Please obtain a copy of the License at
# http://opensource.org/licenses/Apache-2.0/ and read it before using this
# file.
#
# The Original Code and all software distributed under the License are
# distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
# EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
# INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
# Please see the License for the specific language governing rights and
# limitations under the License.
#
# @APPPLANT_LICENSE_HEADER_END@

namespace :test do
  desc 'run integration tests'
  task bintest: [:compile] do
    Go::Build.builds.each do |gb|
      next unless gb.bintest?

      testBinPath = "#{APP_ROOT}/bintest/testFolder/bin/goo"
      bin_path = "#{build_path}/#{gb.name}/bin/goo"
      bin_path << '.exe' if File.exist? "#{bin_path}.exe"
      testBinPath << '.exe' if File.exist? "#{bin_path}.exe"

      sh "cp #{bin_path} #{testBinPath}"
      sh "ruby #{APP_ROOT}/bintest/goo.rb #{testBinPath}"
    end
  end
end
