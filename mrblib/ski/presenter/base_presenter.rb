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
  # Base class for all presenters
  class BasePresenter < BasicObject
    # Initialize a new presenter object.
    #
    # @param [ Hash ] spec The parsed command line arguments.
    #
    # @return [ Void ]
    def initialize(spec)
      @spec = spec
    end

    protected

    # Colorize the output of the task result.
    #
    # @param [ SKI::Result ] result The task result to colorize.
    #
    # @return [ String ]
    def colorize(result)
      colorize_text(result.output, result.success)
    end

    # Colorize the text depend on the given flags.
    #
    # @param [ String ] text      The text to colorize.
    # @param [ Boolean ] no_error Set to false if its an error message.
    #
    # @return [ String ]
    def colorize_text(text, no_error = true)
      if text.nil? || @spec[:'no-color'] || no_error
        text
      else
        text.split("\n").map! { |s| s.set_color(:red) }.join("\n")
      end
    end

    # Adjust the width of the result task output.
    #
    # @param [ SKI::Result ] result The task result to adjust.
    #
    # @return [ String ]
    def adjust(result)
      if @spec[:width] == 0
        result.output
      else
        result.output.scan(/.{1,#{@spec[:width]}}/).join("\n")
      end
    end
  end
end
