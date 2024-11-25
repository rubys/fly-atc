module ApplicationHelper
  def action_cable_meta_tag_dynamic
    scheme = (request.env['HTTP_X_FORWARDED_PROTO'] || request.env["rack.url_scheme"] || '').split(',').last
    return '' if scheme.blank?
    host = request.env['HTTP_X_FORWARDED_HOST'] || request.env["HTTP_HOST"]
    scope = ENV.fetch('FLY_ATC_SCOPE', "")
    root = request.env['RAILS_RELATIVE_URL_ROOT']

    if scope != ""
      websocket = "#{scheme.sub('http', 'ws')}://#{host}#{root}#{scope}/cable"
    else
      websocket = "#{scheme.sub('http', 'ws')}://#{host}#{root}/cable"
    end

    "<meta name=\"action-cable-url\" content=\"#{websocket}\" />".html_safe
  end
end
