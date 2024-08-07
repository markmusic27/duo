import WidgetKit
import SwiftUI


struct Provider: TimelineProvider {
    func placeholder(in context: Context) -> TasksEntry {
        TasksEntry()
    }

    func getSnapshot(in context: Context, completion: @escaping (TasksEntry) -> ()) {
        let entry = TasksEntry()
        completion(entry)
    }

    func getTimeline(in context: Context, completion: @escaping (Timeline<TasksEntry>) -> ()) {
        let entries = [TasksEntry()]  // Single entry
        let timeline = Timeline(entries: entries, policy: .never)
        completion(timeline)
    }
}

struct TasksEntry: TimelineEntry {
    let date = Date()  // Required by TimelineEntry protocol, but not used
}

struct TasksWidgetEntryView: View {
    var entry: Provider.Entry
  
  let exampleTasks: [Task] = [
      Task(title: "Finish differential equations problem set", description: "Pages 12 through 14 in Canvas PDF", dueDate: "Yesterday @ 11:00 PM", priority: 1, isComplete: false),
      Task(title: "Submit tax documents", description: "Gather and submit all tax-related documents", dueDate: "Tomorrow @ 5:00 PM", priority: 2, isComplete: false),
      Task(title: "Buy groceries", description: "Milk, Eggs, Bread, and Butter", dueDate: "Today @ 6:00 PM", priority: 3, isComplete: false),
  ]

    var body: some View {
      ZStack(alignment: .top) {
        Color(hex: 0x191919)
        NavBar()
        
        
//        Text("tasks...")
//          .foregroundColor(.white)
//          .frame(maxWidth: .infinity, maxHeight: .infinity)
//          .multilineTextAlignment(.center)
      }
      .containerBackground(for: .widget) {
        Color(hex: 0x191919)
      }
    }
}

struct CustomCircularProgressViewStyle: ProgressViewStyle {
    var lineWidth: CGFloat

    func makeBody(configuration: Configuration) -> some View {
        ZStack {
            Circle()
                .stroke(Color.white.opacity(0.3), lineWidth: lineWidth)
            Circle()
                .trim(from: 0, to: CGFloat(configuration.fractionCompleted ?? 0))
                .stroke(Color.white.opacity(0.75), style: StrokeStyle(lineWidth: lineWidth, lineCap: .round, lineJoin: .round))
                .rotationEffect(.degrees(-90))
        }
    }
}

struct Task: Identifiable {
  let id = UUID()
  let title: String
  let description: String
  let dueDate: String
  let priority: Int
  let isComplete: Bool
}

struct TasksWidget: Widget {
    let kind: String = "TasksWidget"

    var body: some WidgetConfiguration {
        StaticConfiguration(kind: kind, provider: Provider()) { entry in
            TasksWidgetEntryView(entry: entry)
        }
        .contentMarginsDisabled()
        .configurationDisplayName("Tasks")
        .description("Displays Workspace tasks")
        .supportedFamilies([.systemLarge])
    }
}

struct TasksWidget_Previews: PreviewProvider {
    static var previews: some View {
        TasksWidgetEntryView(entry: TasksEntry())
            .previewContext(WidgetPreviewContext(family: .systemLarge))
    }
}

extension Color {
    init(hex: UInt, alpha: Double = 1) {
        self.init(
            .sRGB,
            red: Double((hex >> 16) & 0xff) / 255,
            green: Double((hex >> 08) & 0xff) / 255,
            blue: Double((hex >> 00) & 0xff) / 255,
            opacity: alpha
        )
    }
}

// Components

struct NavBar: View {
  var body: some View {
    VStack (spacing: 0) {
      ZStack {
        LinearGradient(gradient: Gradient(colors: [Color(hex: 0xFDD600), Color(hex:0xE09400)]),
                                   startPoint: .top,
                                   endPoint: .bottom).overlay(
                                    Rectangle()
                                      .fill(Color.black.opacity(0.2))
                                      .frame(height: 1.6)
                                        .offset(y: 0)
                                    , alignment: .bottom
                                )
        HStack(spacing: 0) {
          ProgressView(value: 0.3)
            .progressViewStyle(CustomCircularProgressViewStyle(lineWidth: 4))
            .frame(height: 18)
            .padding(.trailing, 10)
          Text("Today")
              .foregroundColor(.white)
              .font(.system(size: 15, weight: .bold, design: .rounded))
          Spacer()
          Image("left-arrow")
            .foregroundColor(.white.opacity(0.6))
          Rectangle().frame(width: 10, height: 0)
          Image("right-arrow")
            .foregroundColor(.white.opacity(0.6))
        }.padding(.horizontal, 14)
      }.frame(height: 42)
      Rectangle()
        .fill(
            LinearGradient(
              gradient: Gradient(colors: [Color(hex: 0x191919).opacity(0.8), Color(hex: 0x191919).opacity(0)]),
                startPoint: .top,
                endPoint: .bottom
            )
        )
        .frame(height: 6)
      Line()
         .stroke(style: StrokeStyle(lineWidth: 1, dash: [5]))
         .frame(height: 2)
         .foregroundColor(Color(hex: 0x3C3B40))
      
    }
  }
}

struct Line: Shape {
    func path(in rect: CGRect) -> Path {
        var path = Path()
        path.move(to: CGPoint(x: 0, y: 0))
        path.addLine(to: CGPoint(x: rect.width, y: 0))
        return path
    }
}
