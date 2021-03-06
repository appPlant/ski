# Apache 2.0 License
#
# Copyright (c) 2018 Sebastian Katzer, appPlant GmbH
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

module SKI
  # Save the output under /results/<job> for iss
  class WebPresenter < BasePresenter
    # Initialize a new presenter object.
    #
    # @param [ Hash ]   spec The parsed command line arguments.
    # @param [ String ] cols The name of the columns.
    #
    # @return [ Void ]
    def initialize(spec, cols)
      super(spec)
      @ts   = Time.now.to_i
      @cols = columns(cols || ['OUTPUT'])
    end

    # Format and print the results to $ORBIT_HOME/report/.
    #
    # @param [ Array<SKI::Result> ] results The results to print out.
    #
    # @return [ Void ]
    def present(results)
      File.open("#{make_report_dir}/#{@ts}.skirep", 'w+') do |f|
        f.write render_timestamp_and_columns
        results.each { |res| f.write render_result(res) }
        STDOUT.puts("Written report to: #{f.path}")
      end
    end

    private

    # Render the  header of each report file.
    #
    # @return [ String ]
    def render_timestamp_and_columns
      "#{@ts}\n#{@cols}"
    end

    # Render result into a string that can be written to the IO pipe.
    #
    # @param [ SKI::Result ] res The result to render.
    #
    # @return [ String ]
    def render_result(res)
      planet = res.planet
      output = res.output.gsub!("\n", '\n') || res.output

      %(\n["#{planet.id}","#{planet.name}",#{res.success},#{output}])
    end

    # Convert the columns into tuples of name and type.
    #
    # @param [ String ] The columns to convert.
    #
    # @return [ String ]
    def columns(cols)
      cols.map! do |name|
        case name[-2, 2]
        when '_s' then [name[0...-2], 'string']
        when '_i' then [name[0...-2], 'int']
        when '_f' then [name[0...-2], 'float']
        else           [name,         'string']
        end
      end.inspect
    end

    # Create all parent directories within $ORBIT_HOME
    #
    # @return [ String ] The report sub folder
    def make_report_dir
      rep_dir     = File.join(ENV['ORBIT_HOME'], 'report')
      rep_job_dir = File.join(rep_dir, @spec[:job])

      [rep_dir, rep_job_dir].each do |dir|
        Dir.mkdir(dir) unless Dir.exist? dir
      end

      rep_job_dir
    end
  end
end
